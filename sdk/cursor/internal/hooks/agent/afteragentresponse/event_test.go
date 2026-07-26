package afteragentresponse

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestDecode_AfterAgentResponse(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`)
	if e.EventName() != "afterAgentResponse" {
		t.Fatalf("EventName = %q, want afterAgentResponse", e.EventName())
	}
	if e.ConversationID != "c1" {
		t.Fatalf("ConversationID = %q, want c1", e.ConversationID)
	}
	if e.Text != "done" {
		t.Fatalf("Text = %q, want done", e.Text)
	}
}

func init() {
	register(testCodec)
}
