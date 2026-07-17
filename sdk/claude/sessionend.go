package claude

import (
	"context"
)

// SessionEnd is the SessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

func init() {
	registerDecoder(EventSessionEnd, decodeAs[SessionEnd])
}

// SessionEnd registers an observe-only SessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *Chain {
	registerObserveHandler(c.registerOwner(), fn)
	return c
}
