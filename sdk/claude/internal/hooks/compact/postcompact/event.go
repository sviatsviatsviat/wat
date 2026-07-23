package postcompact

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PostCompact hook event.
type Event struct {
	event.Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PostCompact }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PostCompact, hookkit.EventDecoder[Event](c))
}
