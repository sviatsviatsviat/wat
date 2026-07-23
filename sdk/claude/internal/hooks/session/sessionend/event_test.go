package sessionend

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestDecode(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"session_id":"s","hook_event_name":"SessionEnd","reason":"clear"}`))
	if err != nil {
		t.Fatal(err)
	}
	typed, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if typed.EventName() != event.SessionEnd || typed.Reason != "clear" {
		t.Fatalf("%+v", typed)
	}
}
