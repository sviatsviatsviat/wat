package event

import (
	"encoding/json"
	"testing"
)

func TestToolFields_BindToolInput(t *testing.T) {
	raw := []byte(`{"tool_name":"Bash","tool_input":{"command":"ls"},"tool_use_id":"tu-1"}`)
	var tf ToolFields
	if err := json.Unmarshal(raw, &tf); err != nil {
		t.Fatal(err)
	}
	tf.BindToolInput(raw)
	if tf.ToolName != "Bash" {
		t.Fatalf("ToolName = %q, want Bash", tf.ToolName)
	}
	if tf.ToolUseID != "tu-1" {
		t.Fatalf("ToolUseID = %q, want tu-1", tf.ToolUseID)
	}
	if string(tf.ToolInput.Raw()) != `{"command":"ls"}` {
		t.Fatalf("ToolInput.Raw() = %s", tf.ToolInput.Raw())
	}
}
