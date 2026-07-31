package event

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToolFields_BindToolInput(t *testing.T) {
	raw := []byte(`{"tool_name":"bash","tool_input":{"command":"pwd"}}`)
	var tf ToolFields
	if err := json.Unmarshal(raw, &tf); err != nil {
		t.Fatal(err)
	}
	if err := tf.BindToolInput(raw); err != nil {
		t.Fatal(err)
	}
	if tf.NativeToolName() != "bash" {
		t.Fatalf("NativeToolName() = %q, want bash", tf.NativeToolName())
	}
	if string(tf.Input().Raw()) != `{"command":"pwd"}` {
		t.Fatalf("Input().Raw() = %s", tf.Input().Raw())
	}
}

func TestToolFields_BindToolInput_MissingToolName(t *testing.T) {
	raw := []byte(`{"toolName":"bash","toolArgs":{"command":"pwd"}}`)
	var tf ToolFields
	if err := json.Unmarshal(raw, &tf); err != nil {
		t.Fatal(err)
	}
	err := tf.BindToolInput(raw)
	if err == nil {
		t.Fatal("expected error for camelCase-only tool fields")
	}
	if !strings.Contains(err.Error(), "tool_name is required") {
		t.Fatalf("err = %v", err)
	}
}
