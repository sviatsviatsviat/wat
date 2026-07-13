package hookkit

import (
	"context"
	"fmt"
	"io"
)

// ServeHooks parameterizes the shared hook process serve loop.
type ServeHooks struct {
	Label            string
	HandlerErrorExit int
	Decode           func(raw []byte) (eventName string, event any, err error)
	Lookup           func(eventName string) (fn func(context.Context, any) (any, error), ok bool)
	OnHandlerError   func(eventName string, err error) int
	IsZeroResult     func(result any) bool
	Encode           func(eventName string, result any) (stdout []byte, exitCode int, err error)
}

// ServeLoop reads stdin, dispatches the registered handler, writes stdout, and returns the exit code.
func ServeLoop(ctx context.Context, in io.Reader, out, errw io.Writer, hooks ServeHooks) int {
	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "%s: read stdin: %v\n", hooks.Label, err)
		return hooks.HandlerErrorExit
	}
	eventName, ev, err := hooks.Decode(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "%s: decode: %v\n", hooks.Label, err)
		return hooks.HandlerErrorExit
	}
	fn, ok := hooks.Lookup(eventName)
	if !ok {
		return 0
	}
	result, err := fn(ctx, ev)
	if err != nil {
		_, _ = fmt.Fprintln(errw, err.Error())
		if hooks.OnHandlerError != nil {
			return hooks.OnHandlerError(eventName, err)
		}
		return hooks.HandlerErrorExit
	}
	if result == nil {
		return 0
	}
	if hooks.IsZeroResult != nil && hooks.IsZeroResult(result) {
		return 0
	}
	stdout, code, err := hooks.Encode(eventName, result)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "%s: encode: %v\n", hooks.Label, err)
		return hooks.HandlerErrorExit
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "%s: write stdout: %v\n", hooks.Label, err)
			return hooks.HandlerErrorExit
		}
	}
	return code
}
