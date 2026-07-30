package posttoolbatch

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_PostToolBatch(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"PostToolBatch"}`, event.PostToolBatch)

	t.Run("context", func(t *testing.T) {
		out, code, err := results{}.Context("batch note").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != event.SuccessExit {
			t.Fatalf("exit = %d", code)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		hso, ok := got["hookSpecificOutput"].(map[string]any)
		if !ok || hso["additionalContext"] != "batch note" {
			t.Fatalf("got %s", out)
		}
	})

	t.Run("block", func(t *testing.T) {
		out, code, err := results{}.Block("stop loop").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != event.SuccessExit {
			t.Fatalf("exit = %d", code)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		if got["decision"] != "block" || got["reason"] != "stop loop" {
			t.Fatalf("got %s", out)
		}
	})
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
