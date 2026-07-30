package run

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Serve merges hooks into a local dialect router (same dialect name appends
// handlers), runs one hook dispatch cycle on os.Stdin / os.Stdout / os.Stderr,
// then os.Exit with the resulting code.
//
// Optional --agent and --event flags in os.Args force dialect and event
// selection (skipping payload detect/peek). When a hint disagrees with the
// payload, Serve warns on stderr and continues with the hint.
func Serve(hooks ...Hooks) {
	r := newRouter()
	contribute(r, hooks)
	code := serve(context.Background(), r, os.Stdin, os.Stdout, os.Stderr, parseServeHints(os.Args[1:]))
	os.Exit(code)
}

func contribute(r Registry, hooks []Hooks) {
	for _, h := range hooks {
		if h == nil {
			continue
		}
		h.Contribute(r)
	}
}

func serve(ctx context.Context, router *router, in io.Reader, out io.Writer, errw io.Writer, hints serveHints) int {
	raw, code := readStdin(in, errw)
	if code != 0 {
		return code
	}

	name, d, code := resolveDialect(router, raw, hints, errw)
	if code != 0 {
		return code
	}

	eventName, code := resolveEvent(d, name, raw, hints, errw)
	if code != 0 {
		return code
	}

	handlers := d.HandlersFor(eventName)
	if len(handlers) == 0 {
		return 0
	}

	event, err := d.Codec().DecodeAs(raw, eventName)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", name, err)
		return 1
	}

	acc, code := invokeHandlers(ctx, name, handlers, event, errw)
	if code != 0 {
		return code
	}
	return writeOutput(out, errw, name, acc)
}

func readStdin(in io.Reader, errw io.Writer) (raw []byte, code int) {
	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: read stdin: %v\n", err)
		return nil, 1
	}
	if len(raw) == 0 {
		_, _ = fmt.Fprintln(errw, "run: empty stdin")
		return nil, 1
	}
	return raw, 0
}

func invokeHandlers(ctx context.Context, dialectName string, handlers []hookkit.HookHandler, event hookkit.Event, errw io.Writer) (hookkit.Output, int) {
	var acc hookkit.Output
	for _, h := range handlers {
		outVal, err := h.Invoke(ctx, event)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: handler: %v\n", dialectName, err)
			return nil, 1
		}
		if outVal == nil || outVal.IsZero() {
			continue
		}
		if acc == nil || acc.IsZero() {
			acc = outVal
		} else {
			merged, warnings, err := acc.Merge(outVal)
			if err != nil {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %v\n", dialectName, err)
				return nil, 1
			}
			for _, w := range warnings {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %s\n", dialectName, w)
			}
			acc = merged
		}
		if acc.Stop() {
			break
		}
	}
	return acc, 0
}

func writeOutput(out, errw io.Writer, dialectName string, acc hookkit.Output) int {
	if acc == nil || acc.IsZero() {
		return 0
	}
	stdout, exitCode, err := acc.Encode()
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: encode: %v\n", dialectName, err)
		return 1
	}
	// Claude Code exit 2: host ignores stdout JSON and feeds stderr to the model.
	if dialectName == "claude" && exitCode == 2 {
		if len(stdout) > 0 {
			if _, err := fmt.Fprintln(errw, strings.TrimRight(string(stdout), "\r\n")); err != nil {
				return 1
			}
		}
		return exitCode
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "run: write stdout: %v\n", err)
			return 1
		}
	}
	return exitCode
}
