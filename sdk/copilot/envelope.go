package copilot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PreToolErrorExit is the exit code when a PreToolUse handler returns an error.
// Copilot command hooks fail-closed on non-zero exits other than 2.
const PreToolErrorExit = 1

// HandlerErrorExit is exit code 1 for handler errors under fail-open policy.
const HandlerErrorExit = PreToolErrorExit

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// PermissionRequest it means deny, and for PostToolUseFailure it carries
// additional_context in stdout.
const WarnExit = 2

// Timestamp accepts RFC3339 timestamp strings.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON accepts RFC3339 strings.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	data = bytesTrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("copilot: timestamp: %w", err)
	}
	if s == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("copilot: timestamp: invalid RFC3339: %w", err)
	}
	t.Time = parsed
	return nil
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

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

// Session returns the session identifier.
func (e Envelope) Session() string {
	return e.SessionID
}

// Transcript returns the transcript path.
func (e Envelope) Transcript() string {
	return e.TranscriptPath
}
