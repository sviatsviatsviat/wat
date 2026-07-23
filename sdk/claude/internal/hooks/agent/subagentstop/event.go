package subagentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Event is the SubagentStop hook event.
type Event struct {
	event.Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
	// AgentTranscriptPath is the subagent transcript path when provided.
	AgentTranscriptPath string `json:"agent_transcript_path,omitempty"`
}

// EventName returns the hook event name.
func (Event) EventName() string { return event.SubagentStop }

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.SubagentStop, hookkit.EventDecoder[Event](c))
}
