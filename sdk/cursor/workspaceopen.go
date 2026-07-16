package cursor

import (
	"context"
)

// WorkspaceOpen is the workspaceOpen hook event.
type WorkspaceOpen struct {
	Envelope
}

// EventName returns the canonical hook event name.
func (WorkspaceOpen) EventName() string { return EventWorkspaceOpen }

func init() {
	registerDecoder(EventWorkspaceOpen, decodeAs[WorkspaceOpen])
}

// WorkspaceOpen registers an observe-only workspaceOpen handler.
func (c *Chain) WorkspaceOpen(fn func(context.Context, Hook[WorkspaceOpen]) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}
