package sessionend

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the SessionEnd hook event.
type Event struct {
	event.Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionEnd }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SessionEnd, hookkit.EventDecoder[Event](c))
}
