package claude

// RawEvent holds an unknown or future hook event decoded from the shared envelope only.
type RawEvent struct {
	Envelope
}

// EventName returns the hook event name from the envelope.
func (e RawEvent) EventName() string {
	if e.HookEventName != "" {
		return e.HookEventName
	}
	return "Unknown"
}
