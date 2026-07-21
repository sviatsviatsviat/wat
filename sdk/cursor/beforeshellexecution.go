package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	codec.Register(EventBeforeShellExecution, hookkit.EventDecoder[BeforeShellExecution](codec))
}

// BeforeShellExecution registers a BeforeShellExecution handler on the chain.
func (c *chain) BeforeShellExecution(fn func(context.Context, run.Hook[BeforeShellExecution], PermissionResults) (PermissionOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[BeforeShellExecution]) (PermissionOutput, error) {
		return fn(ctx, hook, permissionResults{})
	}))
	return c
}
