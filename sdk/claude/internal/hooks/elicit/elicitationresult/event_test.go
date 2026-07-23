package elicitationresult

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func TestDecode_ElicitationResult(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"ElicitationResult","server_name":"srv","action":"accept"}`, event.ElicitationResult)
}

func init() {
	Register(runtime.Codec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != wantName {
		t.Fatalf("EventName() = %q, want %q", ev.EventName(), wantName)
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}
