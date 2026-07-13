package cursorhook

import "encoding/json"

// Event is implemented by every decoded Cursor hook event.
type Event interface {
	EventName() string
}

// Attachment is a file or rule attachment on prompt and read-file events.
type Attachment struct {
	// Type is the attachment type.
	Type string `json:"type"`
	// FilePath is the attached file path.
	FilePath string `json:"file_path"`
}

// Edit is a single file edit on afterFileEdit and afterTabFileEdit events.
type Edit struct {
	// OldString is the text replaced.
	OldString string `json:"old_string"`
	// NewString is the replacement text.
	NewString string `json:"new_string"`
}

// SessionStart is the sessionStart hook event.
type SessionStart struct {
	Envelope
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
}

// EventName returns the canonical hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

// BeforeSubmitPrompt is the beforeSubmitPrompt hook event.
type BeforeSubmitPrompt struct {
	Envelope
	// Prompt is the user prompt about to be submitted.
	Prompt string `json:"prompt"`
	// Attachments are context attachments associated with the prompt.
	Attachments []Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (BeforeSubmitPrompt) EventName() string { return EventBeforeSubmitPrompt }

// PreToolUse is the preToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the canonical hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

// ShellCommand extracts the shell command when the tool is Shell.
func (e PreToolUse) ShellCommand() string {
	if e.ToolName != "Shell" && e.ToolName != "shell" {
		return ""
	}
	return extractShellCommand(e.ToolInput)
}

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolOutput is the tool output text.
	ToolOutput string `json:"tool_output"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// DurationMillis returns the execution duration in milliseconds.
func (e PostToolUse) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ErrorMessage is the failure message.
	ErrorMessage string `json:"error_message"`
	// FailureType is the failure category.
	FailureType string `json:"failure_type"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// DurationMillis returns the execution duration in milliseconds.
func (e PostToolUseFailure) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// BeforeShellExecution is the beforeShellExecution hook event.
type BeforeShellExecution struct {
	Envelope
	// Command is the shell command about to run.
	Command string `json:"command"`
	// Sandbox reports whether the command runs in a sandbox.
	Sandbox bool `json:"sandbox"`
}

// EventName returns the canonical hook event name.
func (BeforeShellExecution) EventName() string { return EventBeforeShellExecution }

// AfterShellExecution is the afterShellExecution hook event.
type AfterShellExecution struct {
	Envelope
	// Command is the shell command that ran.
	Command string `json:"command"`
	// Output is the terminal output.
	Output string `json:"output"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (AfterShellExecution) EventName() string { return EventAfterShellExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e AfterShellExecution) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// BeforeMCPExecution is the beforeMCPExecution hook event.
type BeforeMCPExecution struct {
	Envelope
	// ToolName is the MCP tool name (MCP: prefix).
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// URL is the remote MCP server URL when present on the wire.
	URL string `json:"url"`
	// Command is the stdio MCP server command when present on the wire.
	Command string `json:"command"`
}

// EventName returns the canonical hook event name.
func (BeforeMCPExecution) EventName() string { return EventBeforeMCPExecution }

// AfterMCPExecution is the afterMCPExecution hook event.
type AfterMCPExecution struct {
	Envelope
	// ToolName is the MCP tool name (MCP: prefix).
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ResultJSON is the MCP result JSON text.
	ResultJSON string `json:"result_json"`
	// Duration is the execution duration in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (AfterMCPExecution) EventName() string { return EventAfterMCPExecution }

// DurationMillis returns the execution duration in milliseconds.
func (e AfterMCPExecution) DurationMillis() int64 {
	if e.DurationMs != 0 {
		return e.DurationMs
	}
	return e.Duration
}

// BeforeReadFile is the beforeReadFile hook event.
type BeforeReadFile struct {
	Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
	// Attachments are additional file attachments.
	Attachments []Attachment `json:"attachments"`
}

// EventName returns the canonical hook event name.
func (BeforeReadFile) EventName() string { return EventBeforeReadFile }

// AfterFileEdit is the afterFileEdit hook event.
type AfterFileEdit struct {
	Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (AfterFileEdit) EventName() string { return EventAfterFileEdit }

// SubagentStart is the subagentStart hook event.
type SubagentStart struct {
	Envelope
	// SubagentID is the subagent identifier.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
}

// EventName returns the canonical hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// SubagentStop is the subagentStop hook event.
type SubagentStop struct {
	Envelope
	// SubagentID is the subagent identifier.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
	// Summary is the subagent result summary.
	Summary string `json:"summary"`
	// Status is the subagent stop status.
	Status string `json:"status"`
	// LoopCount is the stop-loop iteration count.
	LoopCount int `json:"loop_count"`
	// AgentTranscriptPath is the subagent transcript path when present.
	AgentTranscriptPath *string `json:"agent_transcript_path"`
}

// EventName returns the canonical hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

// Stop is the stop hook event.
type Stop struct {
	Envelope
	// Status is the stop status.
	Status string `json:"status"`
	// LoopCount is the stop-loop iteration count.
	LoopCount int `json:"loop_count"`
}

// EventName returns the canonical hook event name.
func (Stop) EventName() string { return EventStop }

// PreCompact is the preCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
}

// EventName returns the canonical hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// AfterAgentResponse is the afterAgentResponse hook event.
type AfterAgentResponse struct {
	Envelope
	// Text is the agent response text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (AfterAgentResponse) EventName() string { return EventAfterAgentResponse }

// AfterAgentThought is the afterAgentThought hook event.
type AfterAgentThought struct {
	Envelope
	// Text is the agent thought text.
	Text string `json:"text"`
}

// EventName returns the canonical hook event name.
func (AfterAgentThought) EventName() string { return EventAfterAgentThought }

// BeforeTabFileRead is the beforeTabFileRead hook event.
type BeforeTabFileRead struct {
	Envelope
	// FilePath is the file path being read.
	FilePath string `json:"file_path"`
	// Content is the file content.
	Content string `json:"content"`
}

// EventName returns the canonical hook event name.
func (BeforeTabFileRead) EventName() string { return EventBeforeTabFileRead }

// AfterTabFileEdit is the afterTabFileEdit hook event.
type AfterTabFileEdit struct {
	Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (AfterTabFileEdit) EventName() string { return EventAfterTabFileEdit }

// WorkspaceOpen is the workspaceOpen hook event.
type WorkspaceOpen struct {
	Envelope
}

// EventName returns the canonical hook event name.
func (WorkspaceOpen) EventName() string { return EventWorkspaceOpen }

// RawEvent holds an unknown or future event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched JSON payload.
	Raw json.RawMessage
}

// EventName returns the received event name or an empty string.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	return e.receivedName
}

func extractShellCommand(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var v struct {
		Command string `json:"command"`
	}
	if json.Unmarshal(input, &v) != nil {
		return ""
	}
	return v.Command
}
