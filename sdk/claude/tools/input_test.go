package tools

import (
	"encoding/json"
	"testing"
)

func TestAsGrep_CaseInsensitive(t *testing.T) {
	in := NewInput(ToolGrep, json.RawMessage(`{"pattern":"foo","-i":true}`))
	got, ok := in.AsGrep()
	if !ok || !got.CaseInsensitive {
		t.Fatalf("AsGrep = %+v, %v", got, ok)
	}
}

func TestAsAskUserQuestion_MultiSelect(t *testing.T) {
	in := NewInput(ToolAskUserQuestion, json.RawMessage(`{"questions":[{"question":"Pick","header":"H","options":[{"label":"A"}],"multiSelect":true}]}`))
	got, ok := in.AsAskUserQuestion()
	if !ok || len(got.Questions) != 1 || !got.Questions[0].MultiSelect {
		t.Fatalf("AsAskUserQuestion = %+v, %v", got, ok)
	}
}

func TestAsBash(t *testing.T) {
	in := NewInput(ToolBash, json.RawMessage(`{"command":"ls -la","description":"list"}`))
	got, ok := in.AsBash()
	if !ok || got.Command != "ls -la" {
		t.Fatalf("AsBash = %+v, %v", got, ok)
	}
	if _, ok := in.AsWrite(); ok {
		t.Fatal("AsWrite should be false for Bash")
	}
}

func TestAsMCPTool(t *testing.T) {
	in := NewInput("mcp__github__create_issue", json.RawMessage(`{"title":"bug"}`))
	got, ok := in.AsMCPTool()
	if !ok || got.Server != "github" || got.Tool != "create_issue" {
		t.Fatalf("AsMCPTool = %+v, %v", got, ok)
	}
	if _, ok := in.AsBash(); ok {
		t.Fatal("AsBash should be false for MCP tool")
	}
}
