package agnostic

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
)

// Kind is a normalized hook event category shared by all supported agents.
type Kind string

const (
	// KindSessionStart is the normalized category for session start events.
	KindSessionStart Kind = "SessionStart"
	// KindSessionEnd is the normalized category for session end events.
	KindSessionEnd Kind = "SessionEnd"
	// KindUserPrompt is the normalized category for user prompt submission events.
	KindUserPrompt Kind = "UserPrompt"
	// KindPreTool is the normalized category for pre-tool hook events.
	KindPreTool Kind = "PreTool"
	// KindPostTool is the normalized category for successful post-tool events.
	KindPostTool Kind = "PostTool"
	// KindPostToolFailure is the normalized category for failed post-tool events.
	KindPostToolFailure Kind = "PostToolFailure"
	// KindPermissionRequest is the normalized category for permission request events.
	KindPermissionRequest Kind = "PermissionRequest"
	// KindSubagentStart is the normalized category for subagent start events.
	KindSubagentStart Kind = "SubagentStart"
	// KindSubagentStop is the normalized category for subagent stop events.
	KindSubagentStop Kind = "SubagentStop"
	// KindStop is the normalized category for agent stop events.
	KindStop Kind = "Stop"
	// KindPreCompact is the normalized category for pre-compaction events.
	KindPreCompact Kind = "PreCompact"
	// KindNotification is the normalized category for notification events.
	KindNotification Kind = "Notification"
	// KindAgentError is the normalized category for agent runtime error events.
	KindAgentError Kind = "AgentError"
	// KindOther is the normalized category for events with no dedicated mapping.
	KindOther Kind = "Other"
)

// Event is the unified, agent-independent view of a hook invocation.
// Raw always carries the untouched native payload, so nothing is lost by
// normalization; agent-specific handlers can re-decode it with native types.
type Event struct {
	// Agent is the dialect that emitted this hook event.
	Agent Dialect
	// Kind is the normalized event category.
	Kind Kind
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

	// Prompt holds the user prompt text for KindUserPrompt events.
	Prompt string
	// Tool holds tool invocation details for pre/post tool and permission events.
	Tool *ToolCall
	// Result holds tool outcome details for post-tool events.
	Result *ToolResult
	// Subagent holds subagent lifecycle details for subagent start/stop events.
	Subagent *subagent
	// Turn holds turn-end details for KindStop events.
	Turn *turnEnd
	// Compact holds compaction details for KindPreCompact events.
	Compact *compactInfo
	// Note holds notification or error details for KindNotification and KindAgentError events.
	Note *note
	// Life holds session lifecycle details for KindSessionStart and KindSessionEnd events.
	Life *lifecycle
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
	// Error is the failure message for KindPostToolFailure events.
	Error string
	// FailureType is the failure category (error, timeout, permission_denied).
	FailureType string
	// DurationMs is the tool execution duration in milliseconds when available.
	DurationMs int64
}

// subagent describes subagent lifecycle events.
type subagent struct {
	ID             string
	Type           string
	Task           string
	Summary        string
	Status         string
	TranscriptPath string
	LoopCount      int
}

// turnEnd describes agent stop events.
type turnEnd struct {
	Status               string
	LoopCount            int
	StopHookActive       bool
	LastAssistantMessage string
}

// compactInfo describes context compaction events.
type compactInfo struct {
	Trigger            string
	CustomInstructions string
}

// note describes notifications and runtime errors.
type note struct {
	Type        string
	Title       string
	Message     string
	Recoverable *bool
}

// lifecycle describes session start and end events.
type lifecycle struct {
	Source        string
	Reason        string
	Model         string
	InitialPrompt string
	Background    bool
}
