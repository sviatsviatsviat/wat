package hookkit

import "encoding/json"

// EnvelopeAccessor exposes decoded raw JSON from an event envelope.
type EnvelopeAccessor interface {
	DecodedRaw() json.RawMessage
}

// RawBytes returns the untouched JSON for an event when available.
// It prefers accessor.DecodedRaw, then falls back to json.Marshal(ev).
func RawBytes(ev any, accessor EnvelopeAccessor) json.RawMessage {
	if accessor != nil {
		if raw := accessor.DecodedRaw(); len(raw) > 0 {
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
