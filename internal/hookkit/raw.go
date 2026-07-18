package hookkit

import "encoding/json"

// RawPayload holds untouched stdin JSON for a decoded hook event.
// Embed it in agent event structs alongside Envelope.
type RawPayload struct {
	raw json.RawMessage `json:"-"`
}

// SetRaw stores a copy of the original payload bytes.
func (r *RawPayload) SetRaw(b []byte) {
	r.raw = CloneBytes(b)
}

// Raw returns a copy of the stored payload bytes.
func (r RawPayload) Raw() json.RawMessage {
	return CloneRaw(r.raw)
}
