package cursor

import "testing"

func TestNewHookDataAfterShellExecution_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterShellExecution",
		"conversation_id": "conv-1",
		"command": "npm test",
		"output": "ok",
		"duration": 1234,
		"sandbox": false
	}`

	commonData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution([]byte(input), commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}

	assertEqual(t, "afterShellExecution", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "npm test", hookData.Command)
	assertEqual(t, "ok", hookData.Output)
	if hookData.Duration != float32(1234) {
		t.Fatalf("Duration: want 1234, got %v", hookData.Duration)
	}
	if hookData.Sandbox {
		t.Fatal("Sandbox: want false, got true")
	}
}

func TestNewHookDataAfterShellExecution_zeroValueFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterShellExecution",
		"command": "",
		"output": "",
		"duration": 0,
		"sandbox": true
	}`
	commonData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution([]byte(input), commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}

	if hookData.Command != "" || hookData.Output != "" || hookData.Duration != float32(0) || !hookData.Sandbox {
		t.Fatalf("unexpected values: %#v", hookData)
	}
}

func TestNewHookDataAfterShellExecution_decimalDuration(t *testing.T) {
	input := `{
		"hook_event_name": "afterShellExecution",
		"command": "npm test",
		"output": "ok",
		"duration": 2841.805,
		"sandbox": false
	}`

	commonData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution([]byte(input), commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}
	if hookData.Duration != float32(2841.805) {
		t.Fatalf("Duration: want 2841.805, got %v", hookData.Duration)
	}
}

func TestNewHookDataAfterShellExecution_invalidJSON(t *testing.T) {
	_, err := newHookDataAfterShellExecution([]byte(`not json`), hookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
