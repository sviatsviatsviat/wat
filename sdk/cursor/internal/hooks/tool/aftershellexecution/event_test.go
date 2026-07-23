package aftershellexecution

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

func TestDecode_AfterShellExecution(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`)
	if e.Command != "ls" || e.Output != "a\nb" {
		t.Fatalf("event=%+v", e)
	}
}

func init() {
	Register(runtime.Codec)
}
