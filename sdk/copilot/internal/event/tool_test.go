package event

import (
	"encoding/json"
	"testing"
)

func TestToolFields_BindToolInput(t *testing.T) {
	raw := []byte(`{"tool_name":"bash","tool_input":{"command":"pwd"}}`)
	var tf ToolFields
	if err := json.Unmarshal(raw, &tf); err != nil {
		t.Fatal(err)
	}
	tf.BindToolInput(raw)
	if tf.NativeToolName() != "bash" {
		t.Fatalf("NativeToolName() = %q, want bash", tf.NativeToolName())
	}
	if string(tf.Input().Raw()) != `{"command":"pwd"}` {
		t.Fatalf("Input().Raw() = %s", tf.Input().Raw())
	}
}
