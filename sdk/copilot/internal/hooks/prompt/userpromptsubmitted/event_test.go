package userpromptsubmitted

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode_UserPromptSubmitted(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"UserPromptSubmit","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","prompt":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Prompt != "hello" || e.EventName() != event.UserPromptSubmitted {
		t.Fatalf("UserPromptSubmitted=%+v", ev)
	}
}
