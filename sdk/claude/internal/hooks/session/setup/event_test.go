package setup

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

func TestDecode(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"session_id":"s","hook_event_name":"Setup","trigger":"init"}`))
	if err != nil {
		t.Fatal(err)
	}
	typed, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if typed.EventName() != event.Setup || typed.Trigger != "init" {
		t.Fatalf("%+v", typed)
	}
}
