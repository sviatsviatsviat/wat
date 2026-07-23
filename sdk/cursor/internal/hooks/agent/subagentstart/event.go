package subagentstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the subagentStart hook event.
type Event struct {
	event.Envelope
	// SubagentID is the subagent identifier.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStart }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SubagentStart, hookkit.EventDecoder[Event](c))
}
