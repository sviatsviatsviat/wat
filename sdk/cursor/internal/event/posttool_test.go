package event

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

func TestEncode_PostToolUpdatedMCPOutput(t *testing.T) {
	out, code, err := NewPostToolResults().Noop().
		WithUpdatedMCPOutput(map[string]any{"modified": "output"}).
		WithAdditionalContext("extra").
		Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"updated_mcp_tool_output"`) || !strings.Contains(s, `"additional_context":"extra"`) {
		t.Fatalf("bad output: %s", s)
	}
	var parsed map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatal(err)
	}
	mcp, ok := parsed["updated_mcp_tool_output"].(map[string]any)
	if !ok || mcp["modified"] != "output" {
		t.Fatalf("updated_mcp_tool_output = %#v", parsed["updated_mcp_tool_output"])
	}
}

func TestEncode_PostToolZero(t *testing.T) {
	out, code, err := NewPostToolResults().Noop().Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if out != nil {
		t.Fatalf("noop should encode silent stdout, got %s", out)
	}
}

func TestMerge_PostTool_updatedMCPOverwriteWarns(t *testing.T) {
	a := NewPostToolResults().Noop().WithUpdatedMCPOutput("a")
	b := NewPostToolResults().Noop().WithUpdatedMCPOutput("b")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] != hookkit.OverwriteWarning("updatedMCPOutput") {
		t.Fatalf("warnings = %v", warnings)
	}
	out, _, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"updated_mcp_tool_output":"b"`) {
		t.Fatalf("merged = %s", out)
	}
}
