package hookkit

import (
	"encoding/json"
	"fmt"
)

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

// RequireHookEventName peeks hook_event_name and rejects empty stdin or a missing name.
// empty, decodeErr, and nameRequired are the caller's sentinel errors.
func RequireHookEventName(raw []byte, empty, decodeErr, nameRequired error) (string, error) {
	if len(raw) == 0 {
		return "", empty
	}
	name, err := PeekHookEventName(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", decodeErr, err)
	}
	if name == "" {
		return "", nameRequired
	}
	return name, nil
}
