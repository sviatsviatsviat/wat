package workspaceopen

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

func TestDecode_WorkspaceOpen(t *testing.T) {
	mustDecode[Event](t, `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`)
}

func init() {
	Register(runtime.Codec)
}
