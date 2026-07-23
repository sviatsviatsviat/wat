package runtime

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Codec decodes Cursor hook stdin payloads.
var Codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

// DefaultReg is the process-wide default run.Registry.
var DefaultReg = run.GetDefaultRegistry()

func dialectOps() run.DialectOps {
	return run.DialectOps{
		Detect: detectPayload,
		Codec:  Codec,
	}
}

// EnsureDialect attaches this dialect to r when missing.
func EnsureDialect(r *run.Registry) {
	r.EnsureDialect(Dialect, dialectOps())
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
