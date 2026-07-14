package cursor

import (
	"encoding/json"
)

// RawEvent holds an unknown or future event with the full payload preserved.
type RawEvent struct {
	Envelope
	// Raw is the untouched JSON payload.
	Raw json.RawMessage
}

// EventName returns the received event name or an empty string.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	return e.receivedName
}
