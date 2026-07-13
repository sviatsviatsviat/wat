package claudehook

import "encoding/json"

// Event is implemented by every decoded Claude Code hook event.
type Event interface {
	EventName() string
}

// SessionStart is the SessionStart hook event.
type SessionStart struct {
	Envelope
	// Source is the session start source (startup, resume, clear, compact).
	Source string `json:"source"`
	// Model is the model name when provided.
	Model string `json:"model,omitempty"`
	// SessionTitle is the session title when provided.
	SessionTitle string `json:"session_title,omitempty"`
}

// EventName returns the hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// Setup is the Setup hook event.
type Setup struct {
	Envelope
	// Trigger is the setup trigger (init, maintenance).
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (Setup) EventName() string { return EventSetup }

// SessionEnd is the SessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

// UserPromptSubmit is the UserPromptSubmit hook event.
type UserPromptSubmit struct {
	Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (UserPromptSubmit) EventName() string { return EventUserPromptSubmit }

// UserPromptExpansion is the UserPromptExpansion hook event.
type UserPromptExpansion struct {
	Envelope
	// ExpansionType is the expansion kind (slash_command, mcp_prompt).
	ExpansionType string `json:"expansion_type"`
	// CommandName is the slash command name.
	CommandName string `json:"command_name"`
	// CommandArgs is the slash command arguments.
	CommandArgs string `json:"command_args"`
	// CommandSource is the command source.
	CommandSource string `json:"command_source"`
	// Prompt is the expanded prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (UserPromptExpansion) EventName() string { return EventUserPromptExpansion }

// PreToolUse is the PreToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name (matcher field).
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

// PostToolUse is the PostToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolResponse is the tool response JSON.
	ToolResponse json.RawMessage `json:"tool_response"`
}

// EventName returns the hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// PostToolUseFailure is the PostToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Error is the failure message.
	Error string `json:"error"`
}

// EventName returns the hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// PostToolBatch is the PostToolBatch hook event.
type PostToolBatch struct {
	Envelope
}

// EventName returns the hook event name.
func (PostToolBatch) EventName() string { return EventPostToolBatch }

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

// PermissionDenied is the PermissionDenied hook event.
type PermissionDenied struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
}

// EventName returns the hook event name.
func (PermissionDenied) EventName() string { return EventPermissionDenied }

// SubagentStart is the SubagentStart hook event.
type SubagentStart struct {
	Envelope
}

// EventName returns the hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// SubagentStop is the SubagentStop hook event.
type SubagentStop struct {
	Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
	// AgentTranscriptPath is the subagent transcript path when provided.
	AgentTranscriptPath string `json:"agent_transcript_path,omitempty"`
}

// EventName returns the hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

// TaskCreated is the TaskCreated hook event.
type TaskCreated struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCreated) EventName() string { return EventTaskCreated }

// TaskCompleted is the TaskCompleted hook event.
type TaskCompleted struct {
	Envelope
	// Task is the task payload JSON.
	Task json.RawMessage `json:"task"`
}

// EventName returns the hook event name.
func (TaskCompleted) EventName() string { return EventTaskCompleted }

// Stop is the Stop hook event.
type Stop struct {
	Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
}

// EventName returns the hook event name.
func (Stop) EventName() string { return EventStop }

// StopFailure is the StopFailure hook event.
type StopFailure struct {
	Envelope
	// ErrorType is the error category (rate_limit, overloaded, …).
	ErrorType string `json:"error_type"`
	// Message is the error message when provided.
	Message string `json:"message"`
}

// EventName returns the hook event name.
func (StopFailure) EventName() string { return EventStopFailure }

// TeammateIdle is the TeammateIdle hook event.
type TeammateIdle struct {
	Envelope
}

// EventName returns the hook event name.
func (TeammateIdle) EventName() string { return EventTeammateIdle }

// Notification is the Notification hook event.
type Notification struct {
	Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// NotificationType is the notification category.
	NotificationType string `json:"notification_type"`
}

// EventName returns the hook event name.
func (Notification) EventName() string { return EventNotification }

// MessageDisplay is the MessageDisplay hook event.
type MessageDisplay struct {
	Envelope
	// TurnID is the turn identifier.
	TurnID string `json:"turn_id"`
	// MessageID is the message identifier.
	MessageID string `json:"message_id"`
	// Index is the message index in the turn.
	Index int `json:"index"`
	// Final is true when this is the final delta.
	Final bool `json:"final"`
	// Delta is the streamed message delta.
	Delta string `json:"delta"`
}

// EventName returns the hook event name.
func (MessageDisplay) EventName() string { return EventMessageDisplay }

// InstructionsLoaded is the InstructionsLoaded hook event.
type InstructionsLoaded struct {
	Envelope
	// FilePath is the loaded instruction file path.
	FilePath string `json:"file_path"`
	// MemoryType is the memory type (User, Project, Local, Managed).
	MemoryType string `json:"memory_type"`
	// LoadReason is why the file was loaded.
	LoadReason string `json:"load_reason"`
	// Globs are glob patterns when applicable.
	Globs []string `json:"globs,omitempty"`
	// TriggerFilePath is the file that triggered loading.
	TriggerFilePath string `json:"trigger_file_path,omitempty"`
	// ParentFilePath is the parent file path when nested.
	ParentFilePath string `json:"parent_file_path,omitempty"`
}

// EventName returns the hook event name.
func (InstructionsLoaded) EventName() string { return EventInstructionsLoaded }

// ConfigChange is the ConfigChange hook event.
type ConfigChange struct {
	Envelope
	// Source is the config source that changed.
	Source string `json:"source"`
}

// EventName returns the hook event name.
func (ConfigChange) EventName() string { return EventConfigChange }

// CwdChanged is the CwdChanged hook event.
type CwdChanged struct {
	Envelope
	// NewCwd is the new working directory.
	NewCwd string `json:"new_cwd"`
	// OldCwd is the previous working directory.
	OldCwd string `json:"old_cwd,omitempty"`
}

// EventName returns the hook event name.
func (CwdChanged) EventName() string { return EventCwdChanged }

// FileChanged is the FileChanged hook event.
type FileChanged struct {
	Envelope
	// FilePath is the changed file path.
	FilePath string `json:"file_path"`
}

// EventName returns the hook event name.
func (FileChanged) EventName() string { return EventFileChanged }

// WorktreeCreate is the WorktreeCreate hook event.
type WorktreeCreate struct {
	Envelope
}

// EventName returns the hook event name.
func (WorktreeCreate) EventName() string { return EventWorktreeCreate }

// WorktreeRemove is the WorktreeRemove hook event.
type WorktreeRemove struct {
	Envelope
	// WorktreePath is the worktree path being removed.
	WorktreePath string `json:"worktree_path"`
}

// EventName returns the hook event name.
func (WorktreeRemove) EventName() string { return EventWorktreeRemove }

// PreCompact is the PreCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger (manual, auto).
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions.
	CustomInstructions string `json:"custom_instructions"`
}

// EventName returns the hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// PostCompact is the PostCompact hook event.
type PostCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the hook event name.
func (PostCompact) EventName() string { return EventPostCompact }

// Elicitation is the Elicitation hook event.
type Elicitation struct {
	Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Message is the elicitation message.
	Message string `json:"message"`
	// Schema is the requested input schema JSON.
	Schema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Elicitation) EventName() string { return EventElicitation }

// ElicitationResult is the ElicitationResult hook event.
type ElicitationResult struct {
	Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Action is the user action (accept, decline, cancel).
	Action string `json:"action"`
	// Content is the response content JSON.
	Content json.RawMessage `json:"content"`
}

// EventName returns the hook event name.
func (ElicitationResult) EventName() string { return EventElicitationResult }

// RawEvent holds an unknown or future hook event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

// EventName returns the hook event name from the envelope.
func (e RawEvent) EventName() string {
	if e.HookEventName != "" {
		return e.HookEventName
	}
	return "Unknown"
}
