package event

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Timestamp is an RFC3339 timestamp on Copilot hook payloads.
type Timestamp = hookkit.Timestamp

// Envelope holds fields shared by every GitHub Copilot hook event payload.
type Envelope struct {
	// HookEventName is the native hook event name on the wire.
	HookEventName string `json:"hook_event_name"`
	// SessionID is the session identifier.
	SessionID string `json:"session_id"`
	// Timestamp is the event timestamp.
	Timestamp Timestamp `json:"timestamp"`
	// Cwd is the working directory.
	Cwd string `json:"cwd"`
	// TranscriptPath is the conversation transcript path.
	TranscriptPath string `json:"transcript_path"`
}
