package elicitation

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func TestDecode_Elicitation(t *testing.T) {
	mustDecode[Event](t, `{"session_id":"s","hook_event_name":"Elicitation","server_name":"srv","message":"confirm?"}`, event.Elicitation)

	out, code, err := results{}.Accept().WithContent(map[string]any{"answer": "yes"}).Encode()
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
	if hso["hookEventName"] != event.Elicitation {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], event.Elicitation)
	}
	if hso["action"] != "accept" {
		t.Fatalf("action = %v", hso["action"])
	}
	content, ok := hso["content"].(map[string]any)
	if !ok || content["answer"] != "yes" {
		t.Fatalf("content = %v", hso["content"])
	}

	// Results has no Noop(); empty output is the zero-state contract.
	if !(output{}).IsZero() {
		t.Fatal("empty elicitation output should be zero")
	}
}

func init() {
	Register(runtime.Codec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
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
