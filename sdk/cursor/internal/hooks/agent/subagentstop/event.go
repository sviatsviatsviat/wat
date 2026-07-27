package subagentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the subagentStop hook event.
//
// Cursor may send telemetry fields beyond the portable intersection
// (description, duration_ms, message_count, tool_call_count, modified_files).
// Follow-up output is built via StopResults.FollowUp; Cursor only consumes
// followup_message when Status is "completed". Host-side loop caps use the
// hooks.json loop_limit handler option, not an SDK input or output field.
type Event struct {
	event.Envelope
	event.SubagentFields
	// Status is the subagent stop status: "completed", "error", or "aborted".
	Status string `json:"status"`
	// Description is a short description of the subagent's purpose.
	Description string `json:"description"`
	// Summary is the output summary from the subagent.
	Summary string `json:"summary"`
	// DurationMs is the subagent execution time in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	// MessageCount is the number of messages exchanged during the subagent session.
	MessageCount int `json:"message_count"`
	// ToolCallCount is the number of tool calls the subagent made.
	ToolCallCount int `json:"tool_call_count"`
	// LoopCount is how many times a subagentStop follow-up has already triggered
	// for this subagent (starts at 0). Further follow-ups are capped by the
	// hooks.json loop_limit handler option (default 5); that option is install
	// config, not an SDK field.
	LoopCount int `json:"loop_count"`
	// ModifiedFiles lists files the subagent modified.
	ModifiedFiles []string `json:"modified_files"`
	// AgentTranscriptPath is the subagent transcript path when present (null on the wire becomes nil).
	AgentTranscriptPath *string `json:"agent_transcript_path"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SubagentStop, hookkit.EventDecoder[Event](c))
}
