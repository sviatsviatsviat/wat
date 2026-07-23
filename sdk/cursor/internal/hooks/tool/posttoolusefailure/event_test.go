package posttoolusefailure

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestDecode_PostToolUseFailure(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`)
	if e.ErrorMessage != "timeout" || e.FailureType != "timeout" {
		t.Fatalf("event=%+v", e)
	}
}

func init() {
	register(testCodec)
}
