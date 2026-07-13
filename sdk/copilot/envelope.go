package copilot

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Format identifies the Copilot hook wire format.
type Format int

const (
	// FormatUnknown means the payload format could not be determined.
	FormatUnknown Format = iota
	// FormatCamel is the camelCase CLI format (sessionId present).
	FormatCamel
	// FormatVSCode is the VS Code compatible format (hook_event_name present).
	FormatVSCode
)

// PreToolErrorExit is the exit code when a preToolUse handler returns an error.
// Copilot command hooks fail-closed on non-zero exits other than 2.
const PreToolErrorExit = 1

// HandlerErrorExit is exit code 1 for handler errors under fail-open policy.
const HandlerErrorExit = PreToolErrorExit

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// permissionRequest it means deny, and for postToolUseFailure it carries
// additionalContext in stdout.
const WarnExit = 2

// Timestamp accepts ms-epoch numbers (camelCase) and ISO-8601 strings (VS Code).
type Timestamp struct {
	time.Time
}

// UnmarshalJSON accepts ms-epoch numbers or ISO-8601 strings.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	data = bytesTrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		parsed, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return fmt.Errorf("copilot: timestamp: invalid RFC3339: %w", err)
		}
		t.Time = parsed
		return nil
	}
	var ms json.Number
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}
	n, err := ms.Int64()
	if err != nil {
		return fmt.Errorf("copilot: timestamp: invalid epoch: %w", err)
	}
	t.Time = time.UnixMilli(n)
	return nil
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

// Envelope holds fields shared by every GitHub Copilot hook event payload.
type Envelope struct {
	// HookEventName is the native hook event name when present on the wire.
	HookEventName string `json:"hook_event_name"`
	// SessionID is the session identifier (VS Code snake_case).
	SessionID string `json:"session_id"`
	// SessionIDCamel is the session identifier (camelCase CLI).
	SessionIDCamel string `json:"sessionId"`
	// Timestamp is the event timestamp.
	Timestamp Timestamp `json:"timestamp"`
	// Cwd is the working directory.
	Cwd string `json:"cwd"`
	// TranscriptPath is the conversation transcript path (VS Code).
	TranscriptPath string `json:"transcript_path"`
	// TranscriptCamel is the conversation transcript path (camelCase).
	TranscriptCamel string `json:"transcriptPath"`

	receivedName string          `json:"-"`
	canonical    string          `json:"-"`
	decodedRaw   json.RawMessage `json:"-"`
}

// Session returns the session identifier from either wire format.
func (e Envelope) Session() string {
	if e.SessionID != "" {
		return e.SessionID
	}
	return e.SessionIDCamel
}

// Transcript returns the transcript path from either wire format.
func (e Envelope) Transcript() string {
	if e.TranscriptPath != "" {
		return e.TranscriptPath
	}
	return e.TranscriptCamel
}

// ReceivedName returns the hook event name as received on the wire.
func (e Envelope) ReceivedName() string {
	return e.receivedName
}
