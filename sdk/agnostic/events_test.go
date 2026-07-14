package agnostic

import (
	"testing"
)

func TestPreToolEventFrom(t *testing.T) {
	ev := &Event{
		Agent: Claude,
		Kind:  KindPreTool,
		Name:  "PreToolUse",
		Tool:  &ToolCall{Shell: "git push"},
	}
	got, err := PreToolEventFrom(ev)
	if err != nil {
		t.Fatal(err)
	}
	if got.Tool.Shell != "git push" {
		t.Fatalf("Tool.Shell = %q", got.Tool.Shell)
	}

	if _, err := PreToolEventFrom(&Event{Kind: KindPostTool}); err == nil {
		t.Fatal("expected kind mismatch error")
	}
}

func TestAnyEventFrom(t *testing.T) {
	ev := &Event{Kind: KindUserPrompt, Prompt: "hi", Name: "UserPromptSubmit"}
	got := AnyEventFrom(ev)
	if got.Kind != KindUserPrompt || got.Prompt != "hi" {
		t.Fatalf("AnyEvent = %+v", got)
	}
}
