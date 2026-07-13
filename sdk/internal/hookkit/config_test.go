package hookkit

import (
	"encoding/json"
	"testing"
)

func TestParseHandler(t *testing.T) {
	t.Parallel()
	type handler struct {
		Command string `json:"command"`
	}
	got, err := ParseHandler[handler](json.RawMessage(`{"command":"wat run"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Command != "wat run" {
		t.Fatalf("Command = %q", got.Command)
	}
	empty, err := ParseHandler[handler](nil)
	if err != nil {
		t.Fatal(err)
	}
	if empty.Command != "" {
		t.Fatal("empty raw should yield zero value")
	}
}

func TestHandlers(t *testing.T) {
	t.Parallel()
	type handler struct {
		Type string `json:"type"`
	}
	raws, err := Handlers(
		handler{Type: "command"},
		handler{Type: "prompt"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(raws) != 2 {
		t.Fatalf("len = %d", len(raws))
	}
	for i, wantType := range []string{"command", "prompt"} {
		var got handler
		if err := json.Unmarshal(raws[i], &got); err != nil {
			t.Fatalf("raws[%d]: unmarshal: %v", i, err)
		}
		if got.Type != wantType {
			t.Fatalf("raws[%d].Type = %q, want %q", i, got.Type, wantType)
		}
	}
}
