package cursor

import "encoding/json"

// HookHostCursor is the program-argument host token and core.HookContext.HookHost value for Cursor hooks.
const HookHostCursor = "cursor"

// CursorHookRunData is Cursor hook stdin parsed once for subcommands.
// T is the event-specific payload type; use struct{} for common-only hooks with EventSpecific nil.
type CursorHookRunData[T any] struct {
	Common        HookDataCommon
	EventSpecific *T
}

// AfterFileEditEditPair is one edit hunk in an afterFileEdit / afterTabFileEdit payload.
type AfterFileEditEditPair struct {
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}

// AfterFileEditFields is the event-specific JSON shape for afterFileEdit.
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
