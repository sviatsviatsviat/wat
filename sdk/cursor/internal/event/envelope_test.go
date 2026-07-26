package event_test

import (
	"encoding/json"
	"reflect"
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
			{"id": "context", "value": "1m"},
			{"id": "effort", "value": "max"}
		],
		"hook_event_name": "afterAgentResponse",
		"cursor_version": "1.7.2",
		"cwd": "/project"
	}`
	var env event.Envelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatal(err)
	}
	if env.Model != "claude-opus-4-7-thinking-max" {
		t.Fatalf("Model = %q, want claude-opus-4-7-thinking-max", env.Model)
	}
	if env.ModelID != "claude-opus-4-7" {
		t.Fatalf("ModelID = %q, want claude-opus-4-7", env.ModelID)
	}
	wantParams := []event.ModelParam{
		{ID: "thinking", Value: "true"},
		{ID: "context", Value: "1m"},
		{ID: "effort", Value: "max"},
	}
	if !reflect.DeepEqual(env.ModelParams, wantParams) {
		t.Fatalf("ModelParams = %#v, want %#v", env.ModelParams, wantParams)
	}
}

func TestEnvelope_OmittedModelIDAndParams(t *testing.T) {
	const raw = `{
		"conversation_id": "c1",
		"generation_id": "g1",
		"model": "gpt-5",
		"hook_event_name": "stop",
		"cursor_version": "1.7.2"
	}`
	var env event.Envelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatal(err)
	}
	if env.ModelID != "" {
		t.Fatalf("ModelID = %q, want empty when omitted", env.ModelID)
	}
	if env.ModelParams != nil {
		t.Fatalf("ModelParams = %#v, want nil when omitted", env.ModelParams)
	}
	if env.Model != "gpt-5" {
		t.Fatalf("Model = %q, want gpt-5", env.Model)
	}
}
