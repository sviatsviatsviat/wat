package event

// ModelParam is a selected model parameter on Cursor hook payloads
// (for example thinking, context, or effort).
type ModelParam struct {
	// ID is the parameter identifier.
	ID string `json:"id"`
	// Value is the parameter value.
	Value string `json:"value"`
}

// Envelope holds fields shared by every Cursor hook event payload.
//
// It mirrors Cursor's documented common input schema, including optional
// model_id and model_params when the host provides them.
type Envelope struct {
	// ConversationID is the Cursor conversation identifier.
	ConversationID string `json:"conversation_id"`
	// GenerationID is the generation identifier for the current turn.
	GenerationID string `json:"generation_id"`
	// Model is the legacy model slug configured for the composer that triggered
	// the hook, when present on the wire.
	Model string `json:"model"`
	// ModelID is the structured ID for the selected model, when available.
	ModelID string `json:"model_id"`
	// ModelParams lists selected model parameters such as thinking, context, or
	// effort, when present on the wire.
	ModelParams []ModelParam `json:"model_params"`
	// HookEventName is the native hook event name when present on the wire.
	HookEventName string `json:"hook_event_name"`
	// CursorVersion is the Cursor application version.
	CursorVersion string `json:"cursor_version"`
	// WorkspaceRoots lists workspace root paths.
	WorkspaceRoots []string `json:"workspace_roots"`
	// UserEmail is the signed-in user email when present.
	UserEmail *string `json:"user_email"`
	// TranscriptPath is the conversation transcript path when present.
	TranscriptPath *string `json:"transcript_path"`
	// Cwd is the working directory.
	Cwd string `json:"cwd"`
	// SessionID is a fallback session identifier when conversation_id is absent.
	SessionID string `json:"session_id"`
}
