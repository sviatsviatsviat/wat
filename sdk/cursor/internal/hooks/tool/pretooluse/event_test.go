package pretooluse

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

func TestToolInput_AsShell(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(Event)
	input, ok := pre.ToolInput.AsShell()
	if !ok || input.Command != "ls" {
		t.Fatalf("AsShell = %+v, %v", input, ok)
	}
}

func TestDecode_PreToolUse(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`)
	if e.ShellCommand() != "ls" {
		t.Fatalf("ShellCommand=%q", e.ShellCommand())
	}
}

func init() {
	Register(runtime.Codec)
}
