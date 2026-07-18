package cursor

import (
	"context"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	hookkit.RawPayload
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

// OnSessionEnd registers an observe-only sessionEnd handler.
func OnSessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *chain {
	return (&chain{}).SessionEnd(fn)
}

// SessionEnd registers another SessionEnd handler on the chain.
func (c *chain) SessionEnd(fn func(context.Context, Hook[SessionEnd]) error) *chain {
	registerObserveHandler(fn)
	return c
}
