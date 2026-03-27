package cursorcore

// HookHostCursor is the argv host token and core.HookContext.HookHost value for Cursor hooks.
const HookHostCursor = "cursor"

// CursorHookRunData is Cursor hook stdin parsed once for subcommands (e.g. wat run templating).
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
