package cwdchanged

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_CwdChanged(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"CwdChanged","new_cwd":"/new","old_cwd":"/old"}`, event.CwdChanged)
	if ev.NewCwd != "/new" || ev.OldCwd != "/old" {
		t.Fatalf("cwd fields = %+v", ev)
	}
}

func TestEncode_WatchPaths(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		out, code, err := results{}.WatchPaths([]string{"/a/.envrc"}).Encode()
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
		if !ok || len(paths) != 1 || paths[0] != "/a/.envrc" {
			t.Fatalf("watchPaths = %v", hso["watchPaths"])
		}
	})

	t.Run("clearNil", func(t *testing.T) {
		out, code, err := results{}.WatchPaths(nil).Encode()
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
		if !ok || len(paths) != 0 {
			t.Fatalf("watchPaths = %v, want empty array", hso["watchPaths"])
		}
	})

	t.Run("clearEmpty", func(t *testing.T) {
		out, code, err := results{}.WatchPaths([]string{}).Encode()
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
		if !ok || len(paths) != 0 {
			t.Fatalf("watchPaths = %v, want empty array", hso["watchPaths"])
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
