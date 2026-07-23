package posttooluse

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_PostToolUse(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`, event.PostToolUse)

	t.Run("context", func(t *testing.T) {
		out, code, err := results{}.Context("tool finished").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != event.SuccessExit {
			t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		hso, ok := got["hookSpecificOutput"].(map[string]any)
		if !ok {
			t.Fatalf("missing hookSpecificOutput: %s", out)
		}
		if hso["hookEventName"] != event.PostToolUse {
			t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], event.PostToolUse)
		}
		if hso["additionalContext"] != "tool finished" {
			t.Fatalf("additionalContext = %v", hso["additionalContext"])
		}
	})

	t.Run("block", func(t *testing.T) {
		out, code, err := results{}.Block("unsafe output").Encode()
		if err != nil {
			t.Fatal(err)
		}
		if code != event.SuccessExit {
			t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
		}
		var got map[string]any
		if err := json.Unmarshal(out, &got); err != nil {
			t.Fatal(err)
		}
		if got["decision"] != "block" || got["reason"] != "unsafe output" {
			t.Fatalf("block fields = %s", out)
		}
		// Block-only responses put decision/reason at top level; no hookSpecificOutput.
		if _, ok := got["hookSpecificOutput"]; ok {
			t.Fatalf("unexpected hookSpecificOutput: %s", out)
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
