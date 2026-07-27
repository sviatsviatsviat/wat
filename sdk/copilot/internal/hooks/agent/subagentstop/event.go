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
	event.AgentIdentity
	// StopReason is the stop reason.
	StopReason string `json:"stop_reason"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStop }

// Reason returns the stop reason.
func (e Event) Reason() string {
	return e.StopReason
}

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SubagentStop, hookkit.EventDecoder[Event](c))
}
