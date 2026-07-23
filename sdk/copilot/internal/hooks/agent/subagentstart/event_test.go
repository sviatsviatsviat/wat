package subagentstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode_SubagentStart(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"SubagentStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"explore","agent_display_name":"Explore","agent_description":"search codebase"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Name() != "explore" || e.DisplayName() != "Explore" || e.EventName() != event.SubagentStart {
		t.Fatalf("SubagentStart=%+v", ev)
	}
}

func TestEncode_SubagentStartContext(t *testing.T) {
	out, code, err := results{}.Context("hint").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additional_context":"hint"`) {
		t.Fatalf("bad output: %s", out)
	}
}
