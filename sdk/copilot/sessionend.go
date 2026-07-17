package copilot

import (
	"context"
)

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the canonical hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

func init() {
	registerDecoder(EventSessionEnd, decodeAs[SessionEnd])
}

// SessionEnd registers an observe-only SessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *Chain {
	registerObserveHandler(fn)
	return c
}
