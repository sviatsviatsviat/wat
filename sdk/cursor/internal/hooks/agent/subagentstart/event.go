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
	// SubagentType is the subagent type. Official Hooks docs and hooks.json
	// matchers list camelCase values such as generalPurpose; live Cursor
	// payloads may use kebab-case such as general-purpose. Normalize before
	// comparing against matcher-style names.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
	// ParentConversationID is the parent conversation identifier.
	ParentConversationID string `json:"parent_conversation_id"`
	// ToolCallID is the tool call identifier that spawned this subagent.
	ToolCallID string `json:"tool_call_id"`
	// SubagentModel is the model the subagent will use. Live Cursor may send
	// automatic-selection sentinels "", "auto", "default", or "inherit"
	// (case-insensitive after trim) in addition to concrete model IDs. Treat
	// those sentinels as unpinned. When Envelope.Model and SubagentModel are
	// the same concrete ID, the Task is pinned; that equality is not inherit.
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
