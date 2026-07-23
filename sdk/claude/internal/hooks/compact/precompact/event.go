package precompact

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the PreCompact hook event.
type Event struct {
	event.Envelope
	// Trigger is the compaction trigger (manual, auto).
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions.
	CustomInstructions string `json:"custom_instructions"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.PreCompact }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.PreCompact, hookkit.EventDecoder[Event](c))
}
