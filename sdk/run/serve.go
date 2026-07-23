package run

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Serve merges hooks into a local dialect router (same dialect name appends
// handlers), runs one hook dispatch cycle on os.Stdin / os.Stdout / os.Stderr,
// then os.Exit with the resulting code.
func Serve(hooks ...Hooks) {
	r := newRouter()
	for _, h := range hooks {
		if h == nil {
			continue
		}
		h.Contribute(r)
	}
	code := serve(context.Background(), r, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(code)
}

func serve(ctx context.Context, router *router, in io.Reader, out io.Writer, errw io.Writer) int {
	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: read stdin: %v\n", err)
		return 1
	}
	if len(raw) == 0 {
		_, _ = fmt.Fprintln(errw, "run: empty stdin")
		return 1
	}

	name, d, ok := router.detect(raw)
	if !ok || d == nil || d.Codec() == nil {
		_, _ = fmt.Fprintln(errw, "run: unknown dialect")
		return 1
	}

	eventName, err := d.Codec().EventName(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", name, err)
		return 1
	}

	handlers := d.HandlersFor(eventName)
	if len(handlers) == 0 {
		return 0
	}

	event, err := d.Codec().Decode(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", name, err)
		return 1
	}

	var acc hookkit.Output
	for _, h := range handlers {
		outVal, err := h.Invoke(ctx, event)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: handler: %v\n", name, err)
			return 1
		}
		if outVal == nil || outVal.IsZero() {
			continue
		}
		if acc == nil || acc.IsZero() {
			acc = outVal
		} else {
			merged, warnings, err := acc.Merge(outVal)
			if err != nil {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %v\n", name, err)
				return 1
			}
			for _, w := range warnings {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %s\n", name, w)
			}
			acc = merged
		}
		if acc.Stop() {
			break
		}
	}

	if acc == nil || acc.IsZero() {
		return 0
	}

	stdout, exitCode, err := acc.Encode()
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: encode: %v\n", name, err)
		return 1
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "run: write stdout: %v\n", err)
			return 1
		}
	}
	return exitCode
}
