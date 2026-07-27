package event

import (
	"encoding/json"
	"testing"
)

func TestToolFields_BindToolInput(t *testing.T) {
	raw := []byte(`{"tool_name":"Shell","tool_input":{"command":"echo hi"},"tool_use_id":"tu-1"}`)
	var tf ToolFields
	if err := json.Unmarshal(raw, &tf); err != nil {
		t.Fatal(err)
	}
	tf.BindToolInput(raw)
	if tf.ToolName != "Shell" {
		t.Fatalf("ToolName = %q, want Shell", tf.ToolName)
	}
	if tf.ToolUseID != "tu-1" {
		t.Fatalf("ToolUseID = %q, want tu-1", tf.ToolUseID)
	}
	if string(tf.ToolInput.Raw()) != `{"command":"echo hi"}` {
		t.Fatalf("ToolInput.Raw() = %s", tf.ToolInput.Raw())
	}
}
