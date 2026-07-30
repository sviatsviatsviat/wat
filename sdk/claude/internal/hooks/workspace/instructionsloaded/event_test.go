package instructionsloaded

import (
	"context"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_InstructionsLoaded(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"InstructionsLoaded","file_path":"/f","memory_type":"Project","load_reason":"session_start"}`, event.InstructionsLoaded)
}

func TestRegisterHandler_ObserveOnly(t *testing.T) {
	c := hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)
	d := hookkit.NewDialect(c)
	var saw Event
	RegisterHandler(d, func(_ context.Context, hook Event) error {
		saw = hook
		return nil
	})
	handlers := d.HandlersFor(Event{}.EventName())
	if len(handlers) != 1 {
		t.Fatalf("handlers = %d, want 1", len(handlers))
	}
	decoded, err := c.Decode([]byte(`{"session_id":"s","hook_event_name":"InstructionsLoaded","file_path":"/f","memory_type":"Project","load_reason":"session_start"}`))
	if err != nil {
		t.Fatal(err)
	}
	out, err := handlers[0].Invoke(context.Background(), decoded)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("observe-only handler returned %#v", out)
	}
	if saw.FilePath != "/f" || saw.LoadReason != "session_start" {
		t.Fatalf("saw %#v", saw)
	}
	RegisterHandler(d, nil)
	if len(d.HandlersFor(Event{}.EventName())) != 1 {
		t.Fatalf("nil fn should not append handlers")
	}
}

func init() {
	register(testCodec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
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
