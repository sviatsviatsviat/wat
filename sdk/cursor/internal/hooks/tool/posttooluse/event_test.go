package posttooluse

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

func TestDecode_PostToolUse(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`)
	if e.ToolOutput != "contents" || e.DurationMillis() != 100 {
		t.Fatalf("event=%+v", e)
	}
}

func init() {
	register(testCodec)
}
