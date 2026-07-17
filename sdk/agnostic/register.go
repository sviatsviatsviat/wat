package agnostic

import (
	"context"
	"io"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Serve reads a hook payload from in, dispatches registered handlers, writes
// encoded stdout to out, diagnostics to errw, and returns the process exit code.
func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer, opts ...run.Option) int {
	return run.Serve(ctx, in, out, errw, opts...)
}

// WithDialect forces a dialect instead of auto-detection.
func WithDialect(d Dialect) run.Option {
	return run.WithDialect(d.String())
}

// WithEvent supplies the native event name for Copilot camelCase payloads that
// omit hook_event_name.
func WithEvent(name string) run.Option {
	return run.WithEvent(name)
}

// WithGetenv injects environment lookup for Detect.
func WithGetenv(getenv func(string) string) run.Option {
	return run.WithGetenv(getenv)
}
