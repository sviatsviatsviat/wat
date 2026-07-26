package posttoolusefailure

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	e := mustDecode[Event](t, `{
		"hook_event_name":"postToolUseFailure",
		"conversation_id":"c1",
		"cwd":"/project",
		"tool_name":"Shell",
		"tool_input":{"command":"npm test"},
		"tool_use_id":"abc123",
		"error_message":"Command timed out after 30s",
		"failure_type":"timeout",
		"duration":5000,
		"is_interrupt":true
	}`)
	if e.ErrorMessage != "Command timed out after 30s" || e.FailureType != "timeout" {
		t.Fatalf("event=%+v", e)
	}
	if e.ToolUseID != "abc123" || e.DurationMillis() != 5000 {
		t.Fatalf("ids/duration=%+v", e)
	}
	if !e.IsInterrupt {
		t.Fatal("IsInterrupt = false, want true")
	}
	cmd, ok := e.ToolInput.AsShell()
	if !ok || cmd.Command != "npm test" {
		t.Fatalf("ToolInput AsShell=%v ok=%v", cmd, ok)
	}
}

func TestDecode_PostToolUseFailure_durationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"postToolUseFailure",
		"conversation_id":"c1",
		"tool_name":"Shell",
		"error_message":"timeout",
		"failure_type":"timeout",
		"duration_ms":50,
		"is_interrupt":false
	}`)
	if e.DurationMillis() != 50 {
		t.Fatalf("DurationMillis=%d, want 50", e.DurationMillis())
	}
	if e.IsInterrupt {
		t.Fatal("IsInterrupt = true, want false")
	}
}

func init() {
	register(testCodec)
}
