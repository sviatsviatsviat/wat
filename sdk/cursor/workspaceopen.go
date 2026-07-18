package cursor

import (
	"context"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// WorkspaceOpen is the workspaceOpen hook event.
type WorkspaceOpen struct {
	Envelope
	hookkit.RawPayload
}

// EventName returns the canonical hook event name.
func (WorkspaceOpen) EventName() string { return EventWorkspaceOpen }

func init() {
	registerDecoder(EventWorkspaceOpen, decodeAs[WorkspaceOpen])
}

// OnWorkspaceOpen registers an observe-only workspaceOpen handler.
func OnWorkspaceOpen(fn func(context.Context, Hook[WorkspaceOpen]) error) *chain {
	return (&chain{}).WorkspaceOpen(fn)
}

// WorkspaceOpen registers another WorkspaceOpen handler on the chain.
func (c *chain) WorkspaceOpen(fn func(context.Context, Hook[WorkspaceOpen]) error) *chain {
	registerObserveHandler(fn)
	return c
}
