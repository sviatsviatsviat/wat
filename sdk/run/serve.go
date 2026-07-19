package run

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Serve reads a hook payload from in, decodes it once, dispatches registered
// handlers on the default registry, writes merged stdout to out, diagnostics
// to errw, and returns the process exit code.
func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...Option) int {
	return defaultRegistry.serve(ctx, in, out, errw, applyOptions(opts...))
}

// Main runs Serve with os.Stdin, os.Stdout, and os.Stderr, then os.Exit.
func Main(opts ...Option) {
	os.Exit(Serve(context.Background(), os.Stdin, os.Stdout, os.Stderr, opts...))
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

	producers := r.producers(dialect, eventName)
	if len(producers) == 0 {
		return 0
	}

	event, err := ops.Codec.Decode(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", dialect, err)
		return 1
	}

	ctx = WithConfig(ctx, cfg)

	var outputs [][]byte
	exitCode := 0
	for _, p := range producers {
		stdout, code, err := p(ctx, event)
		if err != nil {
			_, _ = fmt.Fprintln(errw, err.Error())
			if code != 0 {
				return code
			}
			return 1
		}
		if len(stdout) > 0 {
			outputs = append(outputs, stdout)
		}
		if code > exitCode {
			exitCode = code
		}
	}

	if len(outputs) == 0 {
		return exitCode
	}

	var merged []byte
	if len(outputs) == 1 {
		merged = outputs[0]
	} else if ops.Merge != nil {
		merged, err = ops.Merge(outputs)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: merge: %v\n", dialect, err)
			return 1
		}
	} else {
		_, _ = fmt.Fprintf(errw, "run: %s: no Merge; discarding %d earlier handler output(s), using last\n", dialect, len(outputs)-1)
		merged = outputs[len(outputs)-1]
	}

	if len(merged) > 0 {
		if _, err := out.Write(merged); err != nil {
			_, _ = fmt.Fprintf(errw, "run: write stdout: %v\n", err)
			return 1
		}
	}
	return exitCode
}
