package erroroccurred

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

func TestErrorOccurred_DetailNull(t *testing.T) {
	e := Event{Error: json.RawMessage("null")}
	if _, ok := e.Detail(); ok {
		t.Fatal("JSON null error payload should be absent")
	}
}

func TestDecode_ErrorOccurred(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"ErrorOccurred","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"error_context":"model_call","recoverable":true}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.EventName() != event.ErrorOccurred {
		t.Fatalf("got %T %+v", ev, ev)
	}
	detail, ok := e.Detail()
	if !ok || detail.Name != "RateLimit" || detail.Message != "slow down" {
		t.Fatalf("Detail=%+v ok=%v", detail, ok)
	}
	if e.Recoverable == nil || !*e.Recoverable {
		t.Fatalf("Recoverable=%v", e.Recoverable)
	}
}
