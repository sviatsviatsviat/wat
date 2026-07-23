package event

import "encoding/json"

// EncodeAdditionalContext renders {"additional_context": ...} Copilot stdout JSON.
func EncodeAdditionalContext(context string) ([]byte, int, error) {
	if context == "" {
		return nil, 0, nil
	}
	out := map[string]any{"additional_context": context}
	b, err := json.Marshal(out)
	return b, 0, err
}
