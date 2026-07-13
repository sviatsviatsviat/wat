package copilot

import (
	"encoding/json"
	"strings"
)

// Event is implemented by every decoded GitHub Copilot hook event.
type Event interface {
	EventName() string
}

// SessionStart is the sessionStart hook event.
type SessionStart struct {
	Envelope
	// Source is the session start source.
	Source string `json:"source"`
	// InitialPrompt is the initial prompt (camelCase).
	InitialPromptCamel string `json:"initialPrompt"`
	// InitialPromptSnake is the initial prompt (VS Code).
	InitialPromptSnake string `json:"initial_prompt"`
}

// EventName returns the canonical hook event name.
func (SessionStart) EventName() string { return EventSessionStart }

// InitialPrompt returns the initial prompt from either wire format.
func (e SessionStart) InitialPrompt() string {
	if e.InitialPromptSnake != "" {
		return e.InitialPromptSnake
	}
	return e.InitialPromptCamel
}

// SessionEnd is the sessionEnd hook event.
type SessionEnd struct {
	Envelope
	// Reason is the session end reason.
	Reason string `json:"reason"`
}

// EventName returns the canonical hook event name.
func (SessionEnd) EventName() string { return EventSessionEnd }

// UserPromptSubmitted is the userPromptSubmitted hook event.
type UserPromptSubmitted struct {
	Envelope
	// Prompt is the submitted user prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the canonical hook event name.
func (UserPromptSubmitted) EventName() string { return EventUserPromptSubmitted }

// PreToolUse is the preToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name (VS Code snake_case).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
}

// EventName returns the canonical hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

// NativeToolName returns the tool name from either wire format.
func (e PreToolUse) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PreToolUse) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e PreToolUse) ShellCommand() string {
	if !isShellToolName(e.NativeToolName()) {
		return ""
	}
	return extractShellCommand(e.Input())
}

// PostToolUse is the postToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
	// ToolResult is the tool result (camelCase).
	ToolResult ToolResult `json:"toolResult"`
	// ToolResultSnake is the tool result (VS Code).
	ToolResultSnake ToolResult `json:"tool_result"`
}

// EventName returns the canonical hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

// NativeToolName returns the tool name from either wire format.
func (e PostToolUse) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PostToolUse) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ResultText returns the textual tool result from either wire format.
func (e PostToolUse) ResultText() string {
	if t := e.ToolResult.Text(); t != "" {
		return t
	}
	return e.ToolResultSnake.Text()
}

// ResultRaw returns the tool result JSON from either wire format.
func (e PostToolUse) ResultRaw() json.RawMessage {
	if raw := extractRawObjectField(e.decodedRawBytes(), "toolResult", "tool_result"); raw != nil {
		return raw
	}
	if e.ToolResult.TextResultForLLM != "" || e.ToolResult.ResultType != "" {
		return marshalToolResultCamel(e.ToolResult)
	}
	if e.ToolResultSnake.TextResultSnake != "" || e.ToolResultSnake.ResultTypeSnake != "" {
		return marshalToolResultSnake(e.ToolResultSnake)
	}
	return nil
}

func extractRawObjectField(raw json.RawMessage, camelKey, snakeKey string) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return nil
	}
	if b, ok := fields[camelKey]; ok && len(b) > 0 && string(b) != "null" {
		return cloneRaw(b)
	}
	if b, ok := fields[snakeKey]; ok && len(b) > 0 && string(b) != "null" {
		return cloneRaw(b)
	}
	return nil
}

func marshalToolResultCamel(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultType != "" {
		out["resultType"] = r.ResultType
	}
	if r.TextResultForLLM != "" {
		out["textResultForLlm"] = r.TextResultForLLM
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil
	}
	return b
}

func marshalToolResultSnake(r ToolResult) json.RawMessage {
	out := map[string]string{}
	if r.ResultTypeSnake != "" {
		out["result_type"] = r.ResultTypeSnake
	}
	if r.TextResultSnake != "" {
		out["text_result_for_llm"] = r.TextResultSnake
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil
	}
	return b
}

// PostToolUseFailure is the postToolUseFailure hook event.
type PostToolUseFailure struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
	// Error is the failure payload (string or object).
	Error json.RawMessage `json:"error"`
}

// EventName returns the canonical hook event name.
func (PostToolUseFailure) EventName() string { return EventPostToolUseFailure }

// NativeToolName returns the tool name from either wire format.
func (e PostToolUseFailure) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PostToolUseFailure) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ErrorMessage returns the failure message from the error field.
func (e PostToolUseFailure) ErrorMessage() string {
	if len(e.Error) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return s
	}
	var detail ErrorDetail
	if json.Unmarshal(e.Error, &detail) == nil {
		return detail.Message
	}
	return string(e.Error)
}

// PermissionRequest is the permissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
}

// EventName returns the canonical hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

// NativeToolName returns the tool name from either wire format.
func (e PermissionRequest) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PermissionRequest) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e PermissionRequest) ShellCommand() string {
	if !isShellToolName(e.NativeToolName()) {
		return ""
	}
	return extractShellCommand(e.Input())
}

// SubagentStart is the subagentStart hook event.
type SubagentStart struct {
	Envelope
	// AgentName is the agent name (VS Code).
	AgentName string `json:"agent_name"`
	// AgentNameCamel is the agent name (camelCase).
	AgentNameCamel string `json:"agentName"`
	// AgentDisplayName is the display name (VS Code).
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDisplayNameCamel is the display name (camelCase).
	AgentDisplayNameCamel string `json:"agentDisplayName"`
	// AgentDescription is the agent description (camelCase).
	AgentDescription string `json:"agentDescription"`
}

// EventName returns the canonical hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// Name returns the agent name from either wire format.
func (e SubagentStart) Name() string {
	if e.AgentName != "" {
		return e.AgentName
	}
	return e.AgentNameCamel
}

// DisplayName returns the agent display name from either wire format.
func (e SubagentStart) DisplayName() string {
	if e.AgentDisplayName != "" {
		return e.AgentDisplayName
	}
	return e.AgentDisplayNameCamel
}

// SubagentStop is the subagentStop hook event.
type SubagentStop struct {
	Envelope
	// AgentName is the agent name (VS Code).
	AgentName string `json:"agent_name"`
	// AgentNameCamel is the agent name (camelCase).
	AgentNameCamel string `json:"agentName"`
	// AgentDisplayName is the display name (VS Code).
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDisplayNameCamel is the display name (camelCase).
	AgentDisplayNameCamel string `json:"agentDisplayName"`
	// StopReason is the stop reason (VS Code).
	StopReason string `json:"stop_reason"`
	// StopReasonCamel is the stop reason (camelCase).
	StopReasonCamel string `json:"stopReason"`
}

// EventName returns the canonical hook event name.
func (SubagentStop) EventName() string { return EventSubagentStop }

// Name returns the agent name from either wire format.
func (e SubagentStop) Name() string {
	if e.AgentName != "" {
		return e.AgentName
	}
	return e.AgentNameCamel
}

// DisplayName returns the agent display name from either wire format.
func (e SubagentStop) DisplayName() string {
	if e.AgentDisplayName != "" {
		return e.AgentDisplayName
	}
	return e.AgentDisplayNameCamel
}

// Reason returns the stop reason from either wire format.
func (e SubagentStop) Reason() string {
	if e.StopReason != "" {
		return e.StopReason
	}
	return e.StopReasonCamel
}

// AgentStop is the agentStop hook event.
type AgentStop struct {
	Envelope
	// StopReason is the stop reason (VS Code).
	StopReason string `json:"stop_reason"`
	// StopReasonCamel is the stop reason (camelCase).
	StopReasonCamel string `json:"stopReason"`
}

// EventName returns the canonical hook event name.
func (AgentStop) EventName() string { return EventAgentStop }

// Reason returns the stop reason from either wire format.
func (e AgentStop) Reason() string {
	if e.StopReason != "" {
		return e.StopReason
	}
	return e.StopReasonCamel
}

// PreCompact is the preCompact hook event.
type PreCompact struct {
	Envelope
	// Trigger is the compaction trigger.
	Trigger string `json:"trigger"`
	// CustomInstructions are user-provided compaction instructions (VS Code).
	CustomInstructions string `json:"custom_instructions"`
	// CustomInstructionsCamel are user-provided compaction instructions (camelCase).
	CustomInstructionsCamel string `json:"customInstructions"`
}

// EventName returns the canonical hook event name.
func (PreCompact) EventName() string { return EventPreCompact }

// Instructions returns custom compaction instructions from either wire format.
func (e PreCompact) Instructions() string {
	if e.CustomInstructions != "" {
		return e.CustomInstructions
	}
	return e.CustomInstructionsCamel
}

// Notification is the notification hook event.
type Notification struct {
	Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// Title is the notification title.
	Title string `json:"title"`
	// NotificationType is the notification category (VS Code).
	NotificationType string `json:"notification_type"`
}

// EventName returns the canonical hook event name.
func (Notification) EventName() string { return EventNotification }

// ErrorOccurred is the errorOccurred hook event.
type ErrorOccurred struct {
	Envelope
	// Error is the error payload (string or object).
	Error json.RawMessage `json:"error"`
	// ErrorContext is additional error context (VS Code).
	ErrorContext string `json:"error_context"`
	// ErrorContextCamel is additional error context (camelCase).
	ErrorContextCamel string `json:"errorContext"`
	// Recoverable is true when the error may be retried.
	Recoverable *bool `json:"recoverable"`
}

// EventName returns the canonical hook event name.
func (ErrorOccurred) EventName() string { return EventErrorOccurred }

// Context returns error context from either wire format.
func (e ErrorOccurred) Context() string {
	if e.ErrorContext != "" {
		return e.ErrorContext
	}
	return e.ErrorContextCamel
}

// Detail parses structured error details when present.
func (e ErrorOccurred) Detail() (ErrorDetail, bool) {
	if len(e.Error) == 0 || string(bytesTrimSpace(e.Error)) == "null" {
		return ErrorDetail{}, false
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return ErrorDetail{Message: s}, true
	}
	var detail ErrorDetail
	if json.Unmarshal(e.Error, &detail) != nil {
		return ErrorDetail{}, false
	}
	return detail, true
}

// ToolResult is a Copilot tool result object in either wire format.
type ToolResult struct {
	// ResultType is the result type (camelCase).
	ResultType string `json:"resultType"`
	// TextResultForLLM is the LLM-facing text (camelCase).
	TextResultForLLM string `json:"textResultForLlm"`
	// ResultTypeSnake is the result type (VS Code).
	ResultTypeSnake string `json:"result_type"`
	// TextResultSnake is the LLM-facing text (VS Code).
	TextResultSnake string `json:"text_result_for_llm"`
}

// Text returns the textual result from either wire format.
func (r ToolResult) Text() string {
	if r.TextResultForLLM != "" {
		return r.TextResultForLLM
	}
	return r.TextResultSnake
}

// ErrorDetail is a structured Copilot error object.
type ErrorDetail struct {
	// Message is the error message.
	Message string `json:"message"`
	// Name is the error category name.
	Name string `json:"name"`
	// Stack is the error stack trace when provided.
	Stack string `json:"stack"`
}

// RawEvent holds an unknown hook event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

// EventName returns the canonical or received hook event name.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	if e.receivedName != "" {
		return e.receivedName
	}
	return "unknown"
}

func extractShellCommand(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var args struct {
		Command string `json:"command"`
	}
	if json.Unmarshal(input, &args) != nil {
		return ""
	}
	return args.Command
}

func isShellToolName(name string) bool {
	switch strings.ToLower(name) {
	case "bash", "powershell", "shell":
		return true
	default:
		return false
	}
}
