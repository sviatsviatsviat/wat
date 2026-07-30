package postcompact

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_PostCompact(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"PostCompact",
		"trigger":"manual",
		"compact_summary":"Summary of the compacted conversation..."
	}`, event.PostCompact)
	if ev.Trigger != "manual" {
		t.Fatalf("Trigger = %q", ev.Trigger)
	}
	if ev.CompactSummary != "Summary of the compacted conversation..." {
		t.Fatalf("CompactSummary = %q", ev.CompactSummary)
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
