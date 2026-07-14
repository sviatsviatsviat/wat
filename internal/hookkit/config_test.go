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

func TestParseFlatCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		raw    string
		want   string
		wantOK bool
	}{
		{name: "command type", raw: `{"type":"command","command":"wat run"}`, want: "wat run", wantOK: true},
		{name: "empty type", raw: `{"command":"wat run"}`, want: "wat run", wantOK: true},
		{name: "prompt type", raw: `{"type":"prompt","command":"ignored"}`, want: "", wantOK: false},
		{name: "invalid json", raw: `{`, want: "", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := ParseFlatCommand(json.RawMessage(tt.raw))
			if ok != tt.wantOK || got != tt.want {
				t.Fatalf("ParseFlatCommand() = (%q, %v), want (%q, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestMarshalFlatCommand_roundTrip(t *testing.T) {
	t.Parallel()
	raw, err := MarshalFlatCommand("wat run --agent copilot --event stop")
	if err != nil {
		t.Fatal(err)
	}
	got, ok := ParseFlatCommand(raw)
	if !ok || got != "wat run --agent copilot --event stop" {
		t.Fatalf("ParseFlatCommand() = (%q, %v)", got, ok)
	}
}
