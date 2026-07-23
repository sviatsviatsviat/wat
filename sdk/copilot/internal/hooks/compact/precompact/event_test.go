package precompact

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

func TestDecode_PreCompact(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","trigger":"auto","custom_instructions":"keep"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.Instructions() != "keep" || e.EventName() != event.PreCompact {
		t.Fatalf("PreCompact=%+v", ev)
	}
}
