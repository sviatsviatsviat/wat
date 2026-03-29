package cursor

import "encoding/json"

// HookHostCursor is the program-argument host token returned by Cursor hook adapters' HookHost method.
const HookHostCursor = "cursor"

// CursorHookRunData is Cursor hook stdin parsed once for subcommands.
// T is the event-specific payload type; use struct{} for common-only hooks with EventSpecific nil.
type CursorHookRunData[T any] struct {
	Common        HookDataCommon
	EventSpecific *T
}

// AfterFileEditEditRange is the range for one edit in an afterTabFileEdit payload (optional on afterFileEdit).
type AfterFileEditEditRange struct {
	StartLineNumber int `json:"start_line_number"`
	StartColumn     int `json:"start_column"`
	EndLineNumber   int `json:"end_line_number"`
	EndColumn       int `json:"end_column"`
}

// AfterFileEditEditPair is one edit hunk in an afterFileEdit / afterTabFileEdit payload.
// afterTabFileEdit may include range, old_line, and new_line; afterFileEdit typically has old_string / new_string only.
type AfterFileEditEditPair struct {
	OldString string                  `json:"old_string"`
	NewString string                  `json:"new_string"`
	EditRange *AfterFileEditEditRange `json:"range,omitempty"`
	OldLine   string                  `json:"old_line,omitempty"`
	NewLine   string                  `json:"new_line,omitempty"`
}

// AfterFileEditFields is the event-specific JSON shape for afterFileEdit and afterTabFileEdit.
type AfterFileEditFields struct {
	FilePath string                  `json:"file_path"`
	Edits    []AfterFileEditEditPair `json:"edits"`
}

// AfterShellExecutionFields is the event-specific JSON shape for afterShellExecution.
type AfterShellExecutionFields struct {
	Command  string  `json:"command"`
	Output   string  `json:"output"`
	Duration float32 `json:"duration"`
	Sandbox  bool    `json:"sandbox"`
}

// AfterMCPExecutionFields is the event-specific JSON shape for afterMCPExecution.
type AfterMCPExecutionFields struct {
	ToolName   string  `json:"tool_name"`
	ToolInput  string  `json:"tool_input"`
	ResultJSON string  `json:"result_json"`
	Duration   float64 `json:"duration"`
}

// AfterAgentResponseFields is the event-specific JSON shape for afterAgentResponse.
type AfterAgentResponseFields struct {
	Text string `json:"text"`
}

// AfterAgentThoughtFields is the event-specific JSON shape for afterAgentThought.
type AfterAgentThoughtFields struct {
	Text       string `json:"text"`
	DurationMs int64  `json:"duration_ms"`
}

// SessionEndFields is the event-specific JSON shape for sessionEnd.
type SessionEndFields struct {
	SessionID         string `json:"session_id"`
	Reason            string `json:"reason"`
	DurationMs        int64  `json:"duration_ms"`
	IsBackgroundAgent bool   `json:"is_background_agent"`
	FinalStatus       string `json:"final_status"`
	ErrorMessage      string `json:"error_message,omitempty"`
}

// HookDataCommon is the shared JSON shape for Cursor hook stdin payloads.
type HookDataCommon struct {
	ConversationID string   `json:"conversation_id"`
	GenerationID   string   `json:"generation_id"`
	Model          string   `json:"model"`
	HookEventName  string   `json:"hook_event_name"`
	CursorVersion  string   `json:"cursor_version"`
	WorkspaceRoots []string `json:"workspace_roots"`
	UserEmail      *string  `json:"user_email"`
	TranscriptPath *string  `json:"transcript_path"`
}

// NewHookDataCommon unmarshals rawJSON into HookDataCommon.
func NewHookDataCommon(rawJSON []byte) (HookDataCommon, error) {
	var hookData HookDataCommon
	if err := json.Unmarshal(rawJSON, &hookData); err != nil {
		return HookDataCommon{}, err
	}
	return hookData, nil
}
