package subagentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode_SubagentStop(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","agent_display_name":"Task","stop_reason":"end_turn"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Name() != "task" || e.Reason() != "end_turn" || e.EventName() != event.SubagentStop {
		t.Fatalf("SubagentStop=%+v", ev)
	}
}

func TestEncode_SubagentStopFollowUp(t *testing.T) {
	out, code, err := agentstop.NewResults().FollowUp("continue").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "continue") {
		t.Fatalf("bad output: %s", out)
	}
}
