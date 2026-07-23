package userpromptsubmitted

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

func TestDecode_UserPromptSubmitted(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"UserPromptSubmit","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","prompt":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Prompt != "hello" || e.EventName() != event.UserPromptSubmitted {
		t.Fatalf("UserPromptSubmitted=%+v", ev)
	}
}
