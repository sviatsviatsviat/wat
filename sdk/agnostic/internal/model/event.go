package model

import (
	"encoding/json"
	"strings"
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
	Subagent *Subagent
	// Turn holds turn-end details for KindStop events.
	Turn *TurnEnd
	// Compact holds compaction details for KindPreCompact events.
	Compact *CompactInfo
	// Note holds notification or error details for KindNotification and KindAgentError events.
	Note *Note
	// Life holds session lifecycle details for KindSessionStart and KindSessionEnd events.
	Life *Lifecycle
}

// ToolCall describes the tool invocation a pre/post tool event refers to.
type ToolCall struct {
	// Name is the normalized tool name (bash, edit, write, read, glob, grep, task, web_fetch, ...).
	Name string
	// Native is the exact original tool name (Bash vs bash vs Shell).
	Native string
	// Input is the native tool input JSON.
	Input json.RawMessage
	// ID is the tool_use_id or tool call id when available.
	ID string
	// Shell is the command string for shell execution tools when available.
	Shell string
	// MCP is true when the call targets an MCP server tool.
	MCP bool
}

// InputAs decodes the native tool input into T.
// A nil ToolCall or empty Input returns the zero value of T with no error.
func InputAs[T any](t *ToolCall) (T, error) {
	var v T
	if t == nil || len(t.Input) == 0 {
		return v, nil
	}
	err := json.Unmarshal(t.Input, &v)
	return v, err
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

// Subagent describes subagent lifecycle events.
type Subagent struct {
	// ID is the subagent identifier.
	ID string
	// Type is the agent type or name.
	Type string
	// Task is the subagent task description.
	Task string
	// Summary is the final output summary when provided.
	Summary string
	// Status is the subagent status (completed, error, aborted, end_turn, ...).
	Status string
	// TranscriptPath is the subagent transcript path when provided.
	TranscriptPath string
	// LoopCount is the number of prior auto follow-ups when reported.
	LoopCount int
}

// TurnEnd describes agent stop events.
type TurnEnd struct {
	// Status is the stop status (completed, aborted, error, end_turn).
	Status string
	// LoopCount is the number of prior auto follow-ups (Cursor).
	LoopCount int
	// StopHookActive is true when a stop hook already forced continuation (Claude).
	StopHookActive bool
	// LastAssistantMessage is the final assistant text of the turn (Claude).
	LastAssistantMessage string
}

// CompactInfo describes context compaction events.
type CompactInfo struct {
	// Trigger is the compaction trigger (manual or auto).
	Trigger string
	// CustomInstructions are user-provided compaction instructions.
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
	// Recoverable is set when the error is recoverable.
	Recoverable *bool
}

// Lifecycle describes session start and end events.
type Lifecycle struct {
	// Source is the session start source (startup, resume, clear, compact, new).
	Source string
	// Reason is the session end reason.
	Reason string
	// Model is the model name when provided at session start.
	Model string
	// InitialPrompt is the initial prompt text when provided at session start.
	InitialPrompt string
	// Background is true for background agent sessions.
	Background bool
}

// Canonical tool names used by ToolCall.Name.
const (
	// ToolBash is the normalized name for shell execution tools.
	ToolBash = "bash"
	// ToolEdit is the normalized name for file edit tools.
	ToolEdit = "edit"
	// ToolWrite is the normalized name for file write tools.
	ToolWrite = "write"
	// ToolRead is the normalized name for file read tools.
	ToolRead = "read"
	// ToolGlob is the normalized name for glob search tools.
	ToolGlob = "glob"
	// ToolGrep is the normalized name for grep search tools.
	ToolGrep = "grep"
	// ToolTask is the normalized name for subagent or task tools.
	ToolTask = "task"
	// ToolWebFetch is the normalized name for web fetch tools.
	ToolWebFetch = "web_fetch"
	// ToolWebSearch is the normalized name for web search tools.
	ToolWebSearch = "web_search"
	// ToolDelete is the normalized name for file delete tools.
	ToolDelete = "delete"
)

var toolAliases = map[string]string{
	// Claude Code
	"bash": ToolBash, "edit": ToolEdit, "write": ToolWrite, "read": ToolRead,
	"glob": ToolGlob, "grep": ToolGrep, "agent": ToolTask, "task": ToolTask,
	"webfetch": ToolWebFetch, "websearch": ToolWebSearch, "notebookedit": ToolEdit,
	// Copilot
	"powershell": ToolBash, "view": ToolRead, "create": ToolWrite, "web_fetch": ToolWebFetch,
	// Cursor
	"shell": ToolBash, "delete": ToolDelete,
}

// NormalizeToolName maps a native tool name onto the canonical vocabulary.
// MCP tools with a verified namespace report mcp=true and keep the native name:
//   - Claude / Copilot PascalCase: "mcp__<server>__<tool>"
//   - Cursor matcher form: "MCP:<tool>"
//
// Copilot camelCase MCP names (serverKey-toolName) are not inferred here; codecs
// set ToolCall.MCP from dialect-specific structured metadata.
func NormalizeToolName(native string) (name string, mcp bool) {
	if strings.HasPrefix(native, "mcp__") || strings.HasPrefix(native, "MCP:") {
		return native, true
	}
	if n, ok := toolAliases[strings.ToLower(native)]; ok {
		return n, false
	}
	return native, false
}
