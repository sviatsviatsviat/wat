package runtime

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Codec decodes Cursor hook stdin payloads.
var Codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

// DefaultDialect is the process-wide Cursor dialect (codec plus handlers).
var DefaultDialect = hookkit.NewDialect(Codec)

// EnsureRegistered attaches DefaultDialect to the process router when missing.
func EnsureRegistered() {
	hookkit.DefaultRouter().Ensure(Dialect, detectPayload, DefaultDialect)
}

func detectPayload(raw []byte, getenv func(string) string) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return true
	}
	if getenv != nil && getenv("CURSOR_VERSION") != "" {
		return true
	}
	return false
}
