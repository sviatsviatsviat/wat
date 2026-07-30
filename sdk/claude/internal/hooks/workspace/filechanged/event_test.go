package filechanged

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_FileChanged(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"FileChanged","file_path":"/Users/my-project/.envrc","event":"change"}`, event.FileChanged)
	if ev.FilePath != "/Users/my-project/.envrc" {
		t.Fatalf("FilePath = %q", ev.FilePath)
	}
	if ev.Change != "change" {
		t.Fatalf("Change = %q", ev.Change)
	}
}

func TestEncode_WatchPaths(t *testing.T) {
	out, code, err := results{}.WatchPaths([]string{"/w/.env"}).Encode()
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
	if !ok {
		t.Fatalf("missing hso: %s", out)
	}
	paths, ok := hso["watchPaths"].([]any)
	if !ok || len(paths) != 1 || paths[0] != "/w/.env" {
		t.Fatalf("watchPaths = %v", hso["watchPaths"])
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
