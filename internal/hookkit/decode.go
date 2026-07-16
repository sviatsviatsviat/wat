package hookkit

import "encoding/json"

// DecodeAsAndThen unmarshals raw into T, then optionally runs after.
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
