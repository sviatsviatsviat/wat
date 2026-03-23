package cursorcore

import "encoding/json"

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
