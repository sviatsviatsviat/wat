package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// AfterShellExecution is the afterShellExecution hook event.
type AfterShellExecution struct {
	Envelope
	// Command is the shell command that ran.
	Command string `json:"command"`
	// Output is the terminal output.
	Output string `json:"output"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (AfterShellExecution) EventName() string { return EventAfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e AfterShellExecution) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

func init() {
	registerDecoder(EventAfterShellExecution, decodeAs[AfterShellExecution])
}

// AfterShellExecution registers an afterShellExecution handler.
func (c *Chain) AfterShellExecution(fn func(context.Context, Hook[AfterShellExecution], PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterShellExecution) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}
