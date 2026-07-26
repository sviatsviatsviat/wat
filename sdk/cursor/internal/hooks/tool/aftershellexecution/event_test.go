package aftershellexecution

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

func TestDecode_AfterShellExecution(t *testing.T) {
	e := mustDecode[Event](t, `{
  "hook_event_name":"afterShellExecution",
  "conversation_id":"c1",
  "command":"ls",
  "output":"a\nb",
  "duration":10,
  "sandbox":true
}`)
	if e.Command != "ls" || e.Output != "a\nb" || e.Duration != 10 || !e.Sandbox {
		t.Fatalf("event=%+v", e)
	}
	if got := e.DurationMillis(); got != 10 {
		t.Fatalf("DurationMillis() = %d, want 10", got)
	}
}

func TestDecode_AfterShellExecution_durationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
  "hook_event_name":"afterShellExecution",
  "conversation_id":"c1",
  "command":"echo hi",
  "output":"hi",
  "duration_ms":25,
  "sandbox":false
}`)
	if e.Sandbox {
		t.Fatalf("Sandbox = true, want false")
	}
	if got := e.DurationMillis(); got != 25 {
		t.Fatalf("DurationMillis() = %d, want 25", got)
	}
}

func init() {
	register(testCodec)
}
