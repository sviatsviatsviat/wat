package cursor

// RawEvent holds an unknown or future hook event decoded from the shared envelope only.
type RawEvent struct {
	Envelope
}

// EventName returns the received event name or an empty string.
func (e RawEvent) EventName() string {
	if e.canonical != "" {
		return e.canonical
	}
	return e.receivedName
}
