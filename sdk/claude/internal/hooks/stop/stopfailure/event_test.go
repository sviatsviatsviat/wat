package stopfailure

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_Failure(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"StopFailure",
		"error":"rate_limit",
		"error_details":"429 Too Many Requests",
		"last_assistant_message":"API Error: Rate limit reached"
	}`, event.StopFailure)
	if ev.Error != "rate_limit" {
		t.Fatalf("Error = %q", ev.Error)
	}
	if ev.ErrorDetails != "429 Too Many Requests" {
		t.Fatalf("ErrorDetails = %q", ev.ErrorDetails)
	}
	if ev.LastAssistantMessage != "API Error: Rate limit reached" {
		t.Fatalf("LastAssistantMessage = %q", ev.LastAssistantMessage)
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
