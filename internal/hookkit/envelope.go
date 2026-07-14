package hookkit

import "encoding/json"

// EnvelopeAccessor exposes decoded raw JSON from an event envelope.
type EnvelopeAccessor interface {
	DecodedRaw() json.RawMessage
}

// RawBytes returns the untouched JSON for an event when available.
func RawBytes(ev any, rawEventRaw json.RawMessage, accessor EnvelopeAccessor, marshal func(any) ([]byte, error)) json.RawMessage {
	if len(rawEventRaw) > 0 {
		return CloneBytes(rawEventRaw)
	}
	if accessor != nil {
		if raw := accessor.DecodedRaw(); len(raw) > 0 {
			return CloneRaw(raw)
		}
	}
	if marshal != nil {
		b, err := marshal(ev)
		if err != nil {
			return nil
		}
		return b
	}
	return nil
}
