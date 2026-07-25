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
	// ParentConversationID is the parent conversation identifier.
	ParentConversationID string `json:"parent_conversation_id"`
	// ToolCallID is the tool call identifier that spawned this subagent.
	ToolCallID string `json:"tool_call_id"`
	// SubagentModel is the model the subagent will use.
	SubagentModel string `json:"subagent_model"`
	// IsParallelWorker reports whether the subagent runs as a parallel worker.
	IsParallelWorker bool `json:"is_parallel_worker"`
	// GitBranch is the git branch the subagent runs on, when present.
	GitBranch string `json:"git_branch"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SubagentStart }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SubagentStart, hookkit.EventDecoder[Event](c))
}
