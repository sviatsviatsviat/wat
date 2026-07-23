package subagentstart

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
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

func TestDecode_SubagentStart(t *testing.T) {
	mustDecode[Event](t, `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`)
}

func init() {
	Register(runtime.Codec)
}
