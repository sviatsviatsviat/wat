package hookkit

import "encoding/json"

// PeekHookEventName extracts hook_event_name from a JSON object payload.
func PeekHookEventName(raw []byte) (string, error) {
	var peek struct {
		HookEventName string `json:"hook_event_name"`
	}
	if err := json.Unmarshal(raw, &peek); err != nil {
		return "", err
	}
	return peek.HookEventName, nil
}
