package hookkit

import "encoding/json"

// RawAccessor exposes untouched JSON from a decoded event.
type RawAccessor interface {
	Raw() json.RawMessage
}

// EventRaw returns the untouched JSON for an event when available.
// It prefers accessor.Raw, then falls back to json.Marshal(ev).
func EventRaw(ev any, accessor RawAccessor) json.RawMessage {
	if accessor != nil {
		if raw := accessor.Raw(); len(raw) > 0 {
			return CloneRaw(raw)
		}
	}
	if ev == nil {
		return nil
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return nil
	}
	return b
}
