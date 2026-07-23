package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the sessionStart hook event.
type Event struct {
	event.Envelope
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionStart }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SessionStart, hookkit.EventDecoder[Event](c))
}
