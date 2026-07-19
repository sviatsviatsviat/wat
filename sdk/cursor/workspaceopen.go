package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// WorkspaceOpen is the workspaceOpen hook event.
type WorkspaceOpen struct {
	Envelope
}

// EventName returns the canonical hook event name.
func (WorkspaceOpen) EventName() string { return EventWorkspaceOpen }

func init() {
	codec.Register(EventWorkspaceOpen, hookkit.EventDecoder[WorkspaceOpen](codec))
}

// WorkspaceOpen registers a WorkspaceOpen handler on the chain.
func (c *chain) WorkspaceOpen(fn func(context.Context, run.Hook[WorkspaceOpen]) error) *chain {
	registerObserveHandler(c.reg, fn)
	return c
}
