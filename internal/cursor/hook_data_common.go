package cursor

import "encoding/json"

// hookDataCommon is the shared JSON shape for Cursor hook stdin payloads.
type hookDataCommon struct {
	ConversationID string   `json:"conversation_id"`
	GenerationID   string   `json:"generation_id"`
	Model          string   `json:"model"`
	HookEventName  string   `json:"hook_event_name"`
	CursorVersion  string   `json:"cursor_version"`
	WorkspaceRoots []string `json:"workspace_roots"`
	UserEmail      *string  `json:"user_email"`
	TranscriptPath *string  `json:"transcript_path"`
}

// newHookDataCommon unmarshals rawJSON into hookDataCommon.
func newHookDataCommon(rawJSON []byte) (hookDataCommon, error) {
	var hookData hookDataCommon
	if err := json.Unmarshal(rawJSON, &hookData); err != nil {
		return hookDataCommon{}, err
	}
	return hookData, nil
}
