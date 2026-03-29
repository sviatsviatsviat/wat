package cursor

import (
	"math"
	"testing"
)

func TestNewHookDataAfterMCPExecution_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterMCPExecution",
		"conversation_id": "conv-1",
		"tool_name": "search",
		"tool_input": "{\"q\":\"x\"}",
		"result_json": "{\"hits\":[]}",
		"duration": 1234
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterMCPExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterMCPExecution", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "search", hookData.Fields.ToolName)
	assertEqual(t, `{"q":"x"}`, hookData.Fields.ToolInput)
	assertEqual(t, `{"hits":[]}`, hookData.Fields.ResultJSON)
	if hookData.Fields.Duration != 1234 {
		t.Fatalf("Duration: want 1234, got %v", hookData.Fields.Duration)
	}
}

func TestNewHookDataAfterMCPExecution_zeroValueFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterMCPExecution",
		"tool_name": "",
		"tool_input": "",
		"result_json": "",
		"duration": 0
	}`
	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterMCPExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	if hookData.Fields.ToolName != "" {
		t.Fatalf("ToolName: want empty string, got %q", hookData.Fields.ToolName)
	}
	if hookData.Fields.ToolInput != "" {
		t.Fatalf("ToolInput: want empty string, got %q", hookData.Fields.ToolInput)
	}
	if hookData.Fields.ResultJSON != "" {
		t.Fatalf("ResultJSON: want empty string, got %q", hookData.Fields.ResultJSON)
	}
	if hookData.Fields.Duration != 0 {
		t.Fatalf("Duration: want 0, got %v", hookData.Fields.Duration)
	}
}

func TestNewHookDataAfterMCPExecution_decimalDuration(t *testing.T) {
	input := `{
		"hook_event_name": "afterMCPExecution",
		"tool_name": "x",
		"tool_input": "{}",
		"result_json": "{}",
		"duration": 2841.805
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterMCPExecutionFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	const epsilon = 1e-6
	wantDuration := 2841.805
	diff := math.Abs(hookData.Fields.Duration - wantDuration)
	if diff > epsilon {
		t.Fatalf("Duration: want ~%v (epsilon %g), got %v", wantDuration, epsilon, hookData.Fields.Duration)
	}
}

func TestNewHookDataAfterMCPExecution_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterMCPExecutionFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
