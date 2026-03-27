package cursor

import (
	"math"
	"testing"
)

func TestNewHookDataAfterShellExecution_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterShellExecution",
		"conversation_id": "conv-1",
		"command": "npm test",
		"output": "ok",
		"duration": 1234,
		"sandbox": false
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterShellExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterShellExecution", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "npm test", hookData.Fields.Command)
	assertEqual(t, "ok", hookData.Fields.Output)
	if hookData.Fields.Duration != float32(1234) {
		t.Fatalf("Duration: want 1234, got %v", hookData.Fields.Duration)
	}
	if hookData.Fields.Sandbox {
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
	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterShellExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	if hookData.Fields.Command != "" {
		t.Fatalf("Command: want empty string, got %q", hookData.Fields.Command)
	}
	if hookData.Fields.Output != "" {
		t.Fatalf("Output: want empty string, got %q", hookData.Fields.Output)
	}
	if hookData.Fields.Duration != float32(0) {
		t.Fatalf("Duration: want 0, got %v", hookData.Fields.Duration)
	}
	if !hookData.Fields.Sandbox {
		t.Fatalf("Sandbox: want true, got false")
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

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterShellExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	const epsilon = 1e-3
	wantDuration := float32(2841.805)
	diff := float64(hookData.Fields.Duration - wantDuration)
	if math.Abs(diff) > epsilon {
		t.Fatalf("Duration: want ~%v (epsilon %g), got %v", wantDuration, epsilon, hookData.Fields.Duration)
	}
}

func TestNewHookDataAfterShellExecution_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterShellExecutionFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
