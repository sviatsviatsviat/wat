package claude

import (
	"encoding/json"
)

// RawEvent holds an unknown or future hook event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched native JSON payload.
	Raw json.RawMessage
}

// EventName returns the hook event name from the envelope.
func (e RawEvent) EventName() string {
	if e.HookEventName != "" {
		return e.HookEventName
	}
	return "Unknown"
}
