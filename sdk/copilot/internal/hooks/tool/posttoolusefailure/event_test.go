package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestEncode_PostToolFailureContext(t *testing.T) {
	out, code, err := results{}.Context("retry with smaller input").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.WarnExit {
		t.Fatalf("code=%d, want %d", code, event.WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestDecode_PostToolUseFailure(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"PostToolUseFailure","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{},"error":"timeout"}`))
	if err != nil {
		t.Fatal(err)
	}
	e := ev.(Event)
	if e.ErrorMessage() != "timeout" {
		t.Fatalf("PostToolUseFailure=%+v", e)
	}
}

func init() {
	register(testCodec)
}
