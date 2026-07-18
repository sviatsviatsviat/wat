package hookkit

import "encoding/json"

// DecodeAsAndThen unmarshals raw into T, then optionally runs after with the
// original stdin bytes (for tool-input construction and similar decode-time work).
func DecodeAsAndThen[T any](raw []byte, after func(*T, []byte)) (T, error) {
	var ev T
	if err := json.Unmarshal(raw, &ev); err != nil {
		return ev, err
	}
	if after != nil {
		after(&ev, raw)
	}
	return ev, nil
}
