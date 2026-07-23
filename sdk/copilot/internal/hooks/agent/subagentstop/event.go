package subagentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
)

// Event is the SubagentStop hook event when the host emits
// hook_event_name "SubagentStop". VS Code-style Stop payloads that include
// agent_name decode as AgentStop instead; use AgentStop.IsSubagent there.
type Event struct {
	event.Envelope
	// AgentName is the agent name.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the display name.
	AgentDisplayName string `json:"agent_display_name"`
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStop }

// Name returns the agent name.
func (e Event) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name.
func (e Event) DisplayName() string {
	return e.AgentDisplayName
}

// Reason returns the stop reason.
func (e Event) Reason() string {
	return e.StopReason
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SubagentStop, hookkit.EventDecoder[Event](c))
}
