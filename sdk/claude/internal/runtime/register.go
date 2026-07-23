package runtime

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Codec decodes Claude Code hook stdin payloads.
var Codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

// DefaultDialect is the process-wide Claude dialect (codec plus handlers).
var DefaultDialect = hookkit.NewDialect(Codec)

// EnsureRegistered attaches DefaultDialect to the process router when missing.
func EnsureRegistered() {
	hookkit.DefaultRouter().Ensure(Dialect, detectPayload, DefaultDialect)
}

func detectPayload(raw []byte) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return false
	}
	if has("hook_event_name") && has("timestamp") {
		return false
	}
	return has("session_id")
}
