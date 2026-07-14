package hookkit

import "encoding/json"

// ToolInputAs decodes raw tool input JSON into T.
func ToolInputAs[T any](raw json.RawMessage) (T, error) {
	var v T
	if len(raw) == 0 {
		return v, nil
	}
	err := json.Unmarshal(raw, &v)
	return v, err
}
