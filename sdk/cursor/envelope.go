package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// HandlerErrorExit is exit code 1. The runner should use this when a handler
// returns an error under Cursor's default fail-open policy.
const HandlerErrorExit = 1

// PermissionDenyExit is exit code 2. Cursor treats it as block/deny on permission-gating
// events, equivalent to returning permission:"deny".
const PermissionDenyExit = 2

// Envelope holds fields shared by every Cursor hook event payload.
type Envelope struct {
	// ConversationID is the Cursor conversation identifier.
	ConversationID string `json:"conversation_id"`
	// GenerationID is the generation identifier for the current turn.
	GenerationID string `json:"generation_id"`
	// Model is the model name when present on the wire.
	Model string `json:"model"`
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

	receivedName string          `json:"-"`
	canonical    string          `json:"-"`
	decodedRaw   json.RawMessage `json:"-"`
}

// Session returns the session identifier from conversation_id or session_id.
func (e Envelope) Session() string {
	if e.ConversationID != "" {
		return e.ConversationID
	}
	return e.SessionID
}

// Transcript returns the transcript path when set.
func (e Envelope) Transcript() string {
	if e.TranscriptPath == nil {
		return ""
	}
	return *e.TranscriptPath
}

// ReceivedName returns the hook event name as received on the wire.
func (e Envelope) ReceivedName() string {
	return e.receivedName
}

func (e *Envelope) setEnvelopeMeta(received, canonical string, raw json.RawMessage) {
	e.receivedName = received
	e.canonical = canonical
	e.decodedRaw = hookkit.CloneRaw(raw)
}

// DecodedRaw returns the untouched JSON stored on the envelope.
func (e Envelope) DecodedRaw() json.RawMessage {
	return hookkit.CloneRaw(e.decodedRaw)
}

// envelope returns a copy of the envelope for Event satisfaction.
// Named envelope (not Envelope) so it does not collide with the embedded field.
func (e Envelope) envelope() Envelope {
	return e
}
