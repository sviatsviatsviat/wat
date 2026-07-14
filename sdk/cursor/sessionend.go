package cursor

import (
	"context"
)

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

func init() {
	registerDecoder(EventSessionEnd, decodeAs[SessionEnd])
}

// SessionEnd registers an observe-only sessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, SessionEndHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}
