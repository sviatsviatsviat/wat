package copilot

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

var codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

var (
	defaultReg   = run.GetDefaultRegistry()
	defaultChain = newChain(defaultReg)
)

func init() {
	ensureDialect(defaultReg)
}

func dialectOps() run.DialectOps {
	return run.DialectOps{
		Detect: detectPayload,
		Codec:  codec,
		Merge:  MergeOutputs,
	}
}

func ensureDialect(r *run.Registry) {
	r.EnsureDialect(Dialect, dialectOps())
}

func detectPayload(raw []byte, getenv func(string) string) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return false
	}
	return has("hook_event_name") && has("timestamp")
}
