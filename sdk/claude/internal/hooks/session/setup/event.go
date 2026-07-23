package setup

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the Setup hook event.
type Event struct {
	event.Envelope
	// Trigger is the setup trigger (init, maintenance).
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.Setup }

// Register registers the Setup decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Setup, hookkit.EventDecoder[Event](c))
}
