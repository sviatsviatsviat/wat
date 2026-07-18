package claude

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// RawEvent holds an unknown or future hook event with the full payload preserved.
type RawEvent struct {
	Envelope
	hookkit.RawPayload
}

// EventName returns the hook event name from the envelope.
func (e RawEvent) EventName() string {
	if e.HookEventName != "" {
		return e.HookEventName
	}
	return "Unknown"
}
