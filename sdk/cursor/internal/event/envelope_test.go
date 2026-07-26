package event_test

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

func TestEnvelope_DecodeModelIDAndParams(t *testing.T) {
	const raw = `{
		"conversation_id": "c1",
		"generation_id": "g1",
		"model": "claude-opus-4-7-thinking-max",
		"model_id": "claude-opus-4-7",
		"model_params": [
			{"id": "thinking", "value": "true"},
			{"id": "context", "value": "1m"}
		],
		"hook_event_name": "stop",
		"cursor_version": "1.7.2"
	}`
	var env event.Envelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatal(err)
	}
	if env.Model != "claude-opus-4-7-thinking-max" {
		t.Fatalf("Model = %q", env.Model)
	}
	if env.ModelID != "claude-opus-4-7" {
		t.Fatalf("ModelID = %q", env.ModelID)
	}
	if len(env.ModelParams) != 2 {
		t.Fatalf("ModelParams len = %d", len(env.ModelParams))
	}
	if env.ModelParams[1].ID != "context" || env.ModelParams[1].Value != "1m" {
		t.Fatalf("ModelParams[1] = %+v", env.ModelParams[1])
	}
}
