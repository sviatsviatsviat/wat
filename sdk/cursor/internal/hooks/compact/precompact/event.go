package precompact

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the preCompact hook event.
type Event struct {
	event.Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PreCompact }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PreCompact, hookkit.EventDecoder[Event](c))
}
