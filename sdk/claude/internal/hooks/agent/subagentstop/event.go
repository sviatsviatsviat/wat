package subagentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the SubagentStop hook event.
type Event struct {
	event.Envelope
	event.StopActiveFields
	// AgentTranscriptPath is the subagent transcript path when provided.
	AgentTranscriptPath string `json:"agent_transcript_path,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.SubagentStop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SubagentStop, hookkit.EventDecoder[Event](c))
}
