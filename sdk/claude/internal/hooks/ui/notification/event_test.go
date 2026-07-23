package notification

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_Notification(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`, event.Notification)

	out, code, err := results{}.Context("idle note").Encode()
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
	if hso["hookEventName"] != event.Notification {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], event.Notification)
	}
	if hso["additionalContext"] != "idle note" {
		t.Fatalf("additionalContext = %v", hso["additionalContext"])
	}
}

func TestMerge_Common_continueFalseStops(t *testing.T) {
	a := results{}.Context("a")
	b := results{}.Context("").WithContinue(false)
	merged, _, err := a.Merge(b)
	if err != nil {
		t.Fatal(err)
	}
	if !merged.Stop() {
		t.Fatal("continue false should stop")
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
