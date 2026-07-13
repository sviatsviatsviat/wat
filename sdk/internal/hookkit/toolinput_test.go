package hookkit

import (
	"encoding/json"
	"testing"
)

func TestToolInputAs(t *testing.T) {
	t.Parallel()
	type input struct {
		Command string `json:"command"`
	}
	got, err := ToolInputAs[input](json.RawMessage(`{"command":"ls"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Command != "ls" {
		t.Fatalf("Command = %q", got.Command)
	}
	empty, err := ToolInputAs[input](nil)
	if err != nil {
		t.Fatal(err)
	}
	if empty.Command != "" {
		t.Fatalf("empty raw Command = %q, want zero value", empty.Command)
	}
}
