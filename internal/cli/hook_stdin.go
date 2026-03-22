package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ReadHookStdinJSON reads stdin, strips a UTF-8 BOM, and validates JSON.
// Whitespace-only input returns (nil, nil). Other non-empty input must be a JSON object.
func ReadHookStdinJSON(stdin io.Reader) ([]byte, error) {
	rawPayload, err := io.ReadAll(stdin)
	if err != nil {
		return nil, err
	}
	rawPayload = bytes.TrimPrefix(rawPayload, []byte{0xEF, 0xBB, 0xBF})

	if len(bytes.TrimSpace(rawPayload)) == 0 {
		return nil, nil
	}

	var jsonObject map[string]json.RawMessage
	if err := json.Unmarshal(rawPayload, &jsonObject); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return rawPayload, nil
}
