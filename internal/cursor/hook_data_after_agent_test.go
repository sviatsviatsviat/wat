package cursor

import "testing"

func TestNewHookDataAfterAgentResponse_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterAgentResponse",
		"conversation_id": "conv-1",
		"text": "Hello from the assistant."
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterAgentResponseFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterAgentResponse", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "Hello from the assistant.", hookData.Fields.Text)
}

func TestNewHookDataAfterAgentThought_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterAgentThought",
		"conversation_id": "conv-2",
		"text": "aggregated thinking",
		"duration_ms": 5000
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterAgentThoughtFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterAgentThought", hookData.HookEventName)
	assertEqual(t, "aggregated thinking", hookData.Fields.Text)
	if hookData.Fields.DurationMs != 5000 {
		t.Fatalf("DurationMs: want 5000, got %d", hookData.Fields.DurationMs)
	}
}

func TestNewHookDataAfterAgentThought_zeroDuration(t *testing.T) {
	input := `{
		"hook_event_name": "afterAgentThought",
		"text": "",
		"duration_ms": 0
	}`
	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterAgentThoughtFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	if hookData.Fields.DurationMs != 0 {
		t.Fatalf("DurationMs: want 0, got %d", hookData.Fields.DurationMs)
	}
}

func TestNewHookDataAfterAgentResponse_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterAgentResponseFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewHookDataAfterAgentThought_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterAgentThoughtFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
