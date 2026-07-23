package run

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Main runs one hook dispatch cycle on os.Stdin / os.Stdout / os.Stderr using
// the default registry, then os.Exit with the resulting code.
func Main(opts ...Option) {
	cfg := applyOptions(opts...)
	os.Exit(GetDefaultRegistry().serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, cfg))
}

func (r *Registry) serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, cfg *Config) int {
	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: read stdin: %v\n", err)
		return 1
	}
	if len(raw) == 0 {
		_, _ = fmt.Fprintln(errw, "run: empty stdin")
		return 1
	}

	dialect, ok := r.detectDialect(raw, cfg.Getenv, cfg.Dialect)
	if !ok {
		_, _ = fmt.Fprintln(errw, "run: unknown dialect")
		return 1
	}

	ops, ok := r.dialectOps(dialect)
	if !ok || ops.Codec == nil {
		_, _ = fmt.Fprintf(errw, "run: %s: missing dialect ops\n", dialect)
		return 1
	}

	eventName, err := ops.Codec.EventName(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", dialect, err)
		return 1
	}

	handlers := r.handlersFor(dialect, eventName)
	if len(handlers) == 0 {
		return 0
	}

	event, err := ops.Codec.Decode(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", dialect, err)
		return 1
	}

	ctx = WithConfig(ctx, cfg)

	var acc Output
	for _, h := range handlers {
		outVal, err := h.invoke(ctx, event)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: handler: %v\n", dialect, err)
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
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %v\n", dialect, err)
				return 1
			}
			for _, w := range warnings {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %s\n", dialect, w)
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
		_, _ = fmt.Fprintf(errw, "run: %s: encode: %v\n", dialect, err)
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
