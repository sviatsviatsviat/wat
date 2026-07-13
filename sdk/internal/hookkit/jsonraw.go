package hookkit

import "encoding/json"

// CloneBytes copies raw bytes into an independent json.RawMessage.
func CloneBytes(raw []byte) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), raw...)
}

// CloneRaw copies raw JSON for independent mutation.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return CloneBytes(raw)
}
