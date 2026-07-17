package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// BeforeShellExecution is the beforeShellExecution hook event.
type BeforeShellExecution struct {
	Envelope
	// Command is the shell command about to run.
	Command string `json:"command"`
	// Sandbox reports whether the command runs in a sandbox.
	Sandbox bool `json:"sandbox"`
}

// EventName returns the canonical hook event name.
func (BeforeShellExecution) EventName() string { return EventBeforeShellExecution }

func init() {
	registerDecoder(EventBeforeShellExecution, decodeAs[BeforeShellExecution])
}

// BeforeShellExecution registers a beforeShellExecution handler.
func (c *Chain) BeforeShellExecution(fn func(context.Context, Hook[BeforeShellExecution], PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeShellExecution) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return c
}
