package model

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// Event is the unified, agent-independent view of a hook invocation.
// Raw always carries the untouched native payload, so nothing is lost by
// normalization; agent-specific handlers can re-decode it with native types.
type Event struct {
	// Agent is the dialect that emitted this hook event (e.g. "claude").
	Agent string
	// Name is the native event name as received (e.g. "beforeShellExecution").
	Name string
	// Session holds session_id, sessionId, or conversation_id from the native payload.
	Session string
	// Cwd is the working directory from the native payload.
	Cwd string
	// TranscriptPath is the conversation transcript path when provided.
	TranscriptPath string
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage

	// Prompt holds the user prompt text for user-prompt events.
	Prompt string
	// Tool holds tool invocation details for pre/post tool and permission events.
	Tool *ToolCall
	// Result holds tool outcome details for post-tool events.
	Result *ToolResult
	// Subagent holds subagent lifecycle details for subagent start/stop events.
	Subagent *Subagent
	// Turn holds turn-end details for stop events.
	Turn *TurnEnd
	// Compact holds compaction details for pre-compaction events.
	Compact *CompactInfo
	// Note holds notification or error details for notification and agent-error events.
	Note *Note
	// Life holds session lifecycle details for session start and end events.
	Life *Lifecycle
}

// ToolCall describes the tool invocation a pre/post tool event refers to.
type ToolCall struct {
	// Name is the normalized tool name (bash, edit, write, read, glob, grep, task, web_fetch, ...).
	Name string
	// Native is the exact original tool name (Bash vs bash vs Shell).
	Native string
	// Input is the typed tool input for this call.
	Input tools.Input
	// ID is the tool_use_id or tool call id when available.
	ID string
	// Shell is the command string for shell execution tools when available.
	Shell string
	// MCP is true when the call targets an MCP server tool.
	MCP bool
}

// ToolResult describes the outcome of a completed or failed tool call.
type ToolResult struct {
	// Text is the textual result as seen by the model, when available.
	Text string
	// Raw is the native result payload JSON.
	Raw json.RawMessage
	// Error is the failure message for post-tool failure events.
	Error string
	// FailureType is the failure category (error, timeout, permission_denied).
	FailureType string
	// DurationMs is the tool execution duration in milliseconds when available.
	DurationMs int64
}

// Subagent describes subagent lifecycle events.
type Subagent struct {
	// ID is the subagent identifier.
	ID string
	// Type is the subagent type.
	Type string
	// Task is the subagent task description.
	Task string
	// Summary is the subagent result summary.
	Summary string
	// Status is the subagent status (completed, error, aborted, end_turn, ...).
	Status string
	// TranscriptPath is the subagent transcript path when provided.
	TranscriptPath string
	// LoopCount is the subagent loop count when available.
	LoopCount int
}

// TurnEnd describes agent stop events.
type TurnEnd struct {
	// Status is the turn-end status when provided.
	Status string
	// LoopCount is the stop loop count when available.
	LoopCount int
	// StopHookActive is true when a stop hook is already active (Claude).
	StopHookActive bool
	// LastAssistantMessage is the last assistant message when provided.
	LastAssistantMessage string
}

// CompactInfo describes context compaction events.
type CompactInfo struct {
	// Trigger is the compaction trigger reason.
	Trigger string
	// CustomInstructions holds extra compaction instructions when provided.
	CustomInstructions string
}

// Note describes notifications and runtime errors.
type Note struct {
	// Type is the notification or error type.
	Type string
	// Title is the notification title when provided.
	Title string
	// Message is the notification or error message.
	Message string
	// Recoverable reports whether an agent error is recoverable when set.
	Recoverable *bool
}

// Lifecycle describes session start and end events.
type Lifecycle struct {
	// Source is the session start source when provided.
	Source string
	// Reason is the session end reason when provided.
	Reason string
	// Model is the model name when provided.
	Model string
	// InitialPrompt is the session start prompt when provided.
	InitialPrompt string
	// Background is true for background agent sessions when provided.
	Background bool
}
