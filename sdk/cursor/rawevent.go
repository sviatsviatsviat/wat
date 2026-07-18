package cursor

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// RawEvent holds an unknown or future event with the full payload preserved.
type RawEvent struct {
	Envelope
	hookkit.RawPayload
}

// EventName returns the received event name or an empty string.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	return e.receivedName
}
