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

// OnSessionEnd registers an observe-only SessionEnd handler.
func OnSessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *chain {
	return (&chain{}).SessionEnd(fn)
}

// SessionEnd registers another SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *chain {
	registerObserveHandler(fn)
	return c
}
