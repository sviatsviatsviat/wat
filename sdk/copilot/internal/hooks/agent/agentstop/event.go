package agentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the Stop (agentStop) hook event.
// When the host scopes the stop to a subagent, AgentName and/or AgentDisplayName
// are set; hook authors should check IsSubagent (or those fields) rather than
// expecting a separate SubagentStop wire name.
type Event struct {
	event.Envelope
	// AgentName is the subagent name when the stop is subagent-scoped.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the subagent display name when present.
	AgentDisplayName string `json:"agent_display_name"`
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AgentStop }

// IsSubagent reports whether this Stop payload is scoped to a subagent.
func (e Event) IsSubagent() bool {
	return e.AgentName != "" || e.AgentDisplayName != ""
}

// Name returns the agent name when present.
func (e Event) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name when present.
func (e Event) DisplayName() string {
	return e.AgentDisplayName
}

// Reason returns the stop reason.
func (e Event) Reason() string {
	return e.StopReason
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AgentStop, hookkit.EventDecoder[Event](c))
}
