package event

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDecisionOutput_Block(t *testing.T) {
	out, code, err := BlockDecision(UserPromptExpansion, "no deploy").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d, want %d", code, SuccessExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["decision"] != "block" || got["reason"] != "no deploy" {
		t.Fatalf("got %s", out)
	}
	if _, ok := got["hookSpecificOutput"]; ok {
		t.Fatalf("unexpected hookSpecificOutput: %s", out)
	}
}

func TestDecisionOutput_Context(t *testing.T) {
	out, code, err := ContextDecision(PreCompact, "keep tests").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d, want %d", code, SuccessExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso, ok := got["hookSpecificOutput"].(map[string]any)
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %s", out)
	}
	if hso["hookEventName"] != PreCompact || hso["additionalContext"] != "keep tests" {
		t.Fatalf("got %s", out)
	}
}

func TestDecisionOutput_MergeBlockWins(t *testing.T) {
	a := ContextDecision(ConfigChange, "note")
	b := BlockDecision(ConfigChange, "blocked")
	merged, _, err := a.Merge(b)
	if err != nil {
		t.Fatal(err)
	}
	out, code, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) {
		t.Fatalf("got %s", out)
	}
	if !merged.Stop() {
		t.Fatal("expected Stop after block")
	}
}

func TestExitBlockOutput_Block(t *testing.T) {
	out, code, err := BlockExitBlock(TeammateIdle, "keep working").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != BlockExit {
		t.Fatalf("exit = %d, want %d", code, BlockExit)
	}
	if string(out) != "keep working" {
		t.Fatalf("stdout body = %q", out)
	}
}

func TestExitBlockOutput_Context(t *testing.T) {
	out, code, err := ContextExitBlock(TaskCreated, "logged").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d, want %d", code, SuccessExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso := got["hookSpecificOutput"].(map[string]any)
	if hso["additionalContext"] != "logged" {
		t.Fatalf("got %s", out)
	}
}

func TestExitBlockOutput_MergeBlockClearsJSON(t *testing.T) {
	a := ContextExitBlock(TaskCompleted, "ctx")
	b := BlockExitBlock(TaskCompleted, "not done")
	merged, _, err := a.Merge(b)
	if err != nil {
		t.Fatal(err)
	}
	out, code, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != BlockExit {
		t.Fatalf("exit = %d, want %d", code, BlockExit)
	}
	if string(out) != "not done" {
		t.Fatalf("body = %q", out)
	}
}
