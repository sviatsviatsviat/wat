package hookkit

import "encoding/json"

// ParseHandler decodes native handler JSON into T.
func ParseHandler[T any](raw json.RawMessage) (T, error) {
	var h T
	if len(raw) == 0 {
		return h, nil
	}
	err := json.Unmarshal(raw, &h)
	return h, err
}

// MarshalHandler encodes h as native handler JSON.
func MarshalHandler[T any](h T) (json.RawMessage, error) {
	return json.Marshal(h)
}

// Handlers encodes typed handlers as native handler JSON blobs.
func Handlers[T any](h ...T) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, 0, len(h))
	for _, handler := range h {
		raw, err := MarshalHandler(handler)
		if err != nil {
			return nil, err
		}
		out = append(out, raw)
	}
	return out, nil
}
