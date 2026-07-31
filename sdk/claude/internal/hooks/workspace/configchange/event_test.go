package configchange

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_ConfigChange(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"ConfigChange","source":"user_settings","file_path":"/Users/me/.claude/settings.json"}`, event.ConfigChange)
	if ev.Source != "user_settings" {
		t.Fatalf("Source = %q", ev.Source)
	}
	if ev.FilePath != "/Users/me/.claude/settings.json" {
		t.Fatalf("FilePath = %q", ev.FilePath)
	}
}

func TestEncode_ConfigChangeBlock(t *testing.T) {
	out, code, err := results{}.Block("reject settings").Encode()
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
	if got["decision"] != "block" || got["reason"] != "reject settings" {
		t.Fatalf("got %s", out)
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
