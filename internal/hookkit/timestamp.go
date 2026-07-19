package hookkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp accepts RFC3339 timestamp strings in JSON.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON accepts RFC3339 strings.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("hookkit: timestamp: %w", err)
	}
	if s == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("hookkit: timestamp: invalid RFC3339: %w", err)
	}
	t.Time = parsed
	return nil
}
