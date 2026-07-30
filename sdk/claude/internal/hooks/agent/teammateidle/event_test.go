package teammateidle

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_TeammateIdle(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"TeammateIdle"}`, event.TeammateIdle)

	out, code, err := results{}.Context("").WithContinue(false).WithStopReason("done").Encode()
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
	if got["continue"] != false || got["stopReason"] != "done" {
		t.Fatalf("got %s", out)
	}
}

func TestEncode_TeammateIdleBlock(t *testing.T) {
	out, code, err := results{}.Block("keep working").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.BlockExit {
		t.Fatalf("exit = %d, want %d", code, event.BlockExit)
	}
	if string(out) != "keep working" {
		t.Fatalf("body = %q", out)
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
