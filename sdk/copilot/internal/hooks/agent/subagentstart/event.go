package subagentstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the SubagentStart hook event.
type Event struct {
	event.Envelope
	// AgentName is the agent name.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the display name.
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDescription is the agent description.
	AgentDescription string `json:"agent_description"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStart }

// Name returns the agent name.
func (e Event) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name.
func (e Event) DisplayName() string {
	return e.AgentDisplayName
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SubagentStart, hookkit.EventDecoder[Event](c))
}
