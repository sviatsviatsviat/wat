package afteragentthought

import (
	"os"
	"path/filepath"
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

func TestDecode_AfterAgentThought(t *testing.T) {
	t.Run("fixture with duration_ms", func(t *testing.T) {
		raw, err := os.ReadFile(filepath.Join("testdata", "after_agent_thought.json"))
		if err != nil {
			t.Fatal(err)
		}
		ev := mustDecode[Event](t, string(raw))
		if ev.EventName() != "afterAgentThought" {
			t.Errorf("EventName = %q, want afterAgentThought", ev.EventName())
		}
		if ev.ConversationID != "c1" {
			t.Errorf("ConversationID = %q, want c1", ev.ConversationID)
		}
		if ev.Text != "fully aggregated thinking text" {
			t.Errorf("Text = %q, want fully aggregated thinking text", ev.Text)
		}
		if ev.DurationMs != 5000 {
			t.Errorf("DurationMs = %d, want 5000", ev.DurationMs)
		}
	})

	t.Run("optional duration_ms omitted", func(t *testing.T) {
		ev := mustDecode[Event](t, `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`)
		if ev.Text != "thinking" {
			t.Errorf("Text = %q, want thinking", ev.Text)
		}
		if ev.DurationMs != 0 {
			t.Errorf("DurationMs = %d, want 0 when omitted", ev.DurationMs)
		}
	})
}

func init() {
	register(testCodec)
}
