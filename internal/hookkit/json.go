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

// RawObjectField returns a copy of the JSON object field named key.
// Missing, empty, and null fields yield nil.
func RawObjectField(raw json.RawMessage, key string) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return nil
	}
	b, ok := fields[key]
	if !ok || len(b) == 0 || string(b) == "null" {
		return nil
	}
	return CloneBytes(b)
}

// RawToText extracts a best-effort textual form of a tool_response value.
func RawToText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}
