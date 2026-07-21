package copilot

import (
	"encoding/json"
	"testing"
)

func TestErrorOccurred_DetailNull(t *testing.T) {
	e := ErrorOccurred{Error: json.RawMessage("null")}
	if _, ok := e.Detail(); ok {
		t.Fatal("JSON null error payload should be absent")
	}
}

func TestDecode_ErrorOccurred(t *testing.T) {
	e := mustDecode[ErrorOccurred](t, `{"hook_event_name":"ErrorOccurred","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"error_context":"model_call","recoverable":true}`, EventErrorOccurred)
	detail, ok := e.Detail()
	if !ok || detail.Name != "RateLimit" || detail.Message != "slow down" {
		t.Fatalf("Detail=%+v ok=%v", detail, ok)
	}
	if e.Recoverable == nil || !*e.Recoverable {
		t.Fatalf("Recoverable=%v", e.Recoverable)
	}
}
