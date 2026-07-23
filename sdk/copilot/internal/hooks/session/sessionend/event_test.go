package sessionend

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

func TestDecode_SessionEnd(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"SessionEnd","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","reason":"complete"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Reason != "complete" || e.EventName() != event.SessionEnd {
		t.Fatalf("SessionEnd=%+v", ev)
	}
}
