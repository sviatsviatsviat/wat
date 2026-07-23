package permissiondenied

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_PermissionDenied(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"PermissionDenied","tool_name":"Bash"}`, event.PermissionDenied)

	if !(output{}).IsZero() {
		t.Fatal("empty permission denied output should be zero")
	}
	retry := results{}.Retry()
	if retry.IsZero() {
		t.Fatal("Retry() should be non-zero")
	}

	out, code, err := retry.Encode()
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
	if hso["hookEventName"] != event.PermissionDenied {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], event.PermissionDenied)
	}
	if hso["retry"] != true {
		t.Fatalf("retry = %v", hso["retry"])
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
