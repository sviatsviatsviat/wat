package claude

import (
	"maps"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEncode_PreToolDeny(t *testing.T) {
	out, code, err := preToolUseResults{}.Deny("destructive command").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out := preToolUseResults{}.Noop()
	if !out.IsZero() {
		t.Fatal("noop should be zero")
	}
}

const preToolUsePayload = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "PreToolUse",
  "tool_name": "Bash",
  "tool_use_id": "tu_1",
  "tool_input": {"command": "rm -rf /tmp/build", "description": "clean"}
}`

func TestDecode_PreToolUse(t *testing.T) {
	ev, err := codec.Decode([]byte(preToolUsePayload))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(PreToolUse)
	if !ok || pre.ToolName != "Bash" || pre.SessionID != "abc123" {
		t.Fatalf("bad event: %+v", ev)
	}
}

func TestMerge_PreToolUse_denyBeatsAllowAndStops(t *testing.T) {
	a := preToolUseResults{}.Allow().WithAdditionalContext("ctx-a")
	b := preToolUseResults{}.Deny("nope").WithAdditionalContext("ctx-b")
	merged, warnings, err := a.Merge(b.(run.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(preToolUseOutput)
	if out.decision != DecisionDeny || out.reason != "nope" {
		t.Fatalf("decision = %q reason = %q", out.decision, out.reason)
	}
	if out.additionalContext != "ctx-a\n\nctx-b" {
		t.Fatalf("context = %q", out.additionalContext)
	}
	if !merged.Stop() {
		t.Fatal("deny should stop")
	}
}

func TestMerge_PreToolUse_updatedInputOverwriteWarns(t *testing.T) {
	first := map[string]any{"command": "echo a"}
	second := map[string]any{"command": "echo b"}
	origFirst := maps.Clone(first)
	a := preToolUseResults{}.Allow().WithUpdatedInput(first)
	b := preToolUseResults{}.Allow().WithUpdatedInput(second)
	merged, warnings, err := a.Merge(b.(run.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] != "updatedInput: overwritten by later handler" {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(preToolUseOutput)
	if out.updatedInput["command"] != "echo b" {
		t.Fatalf("updatedInput = %v", out.updatedInput)
	}
	if !maps.Equal(first, origFirst) {
		t.Fatalf("caller map mutated: %v", first)
	}
}

func TestMerge_PreToolUse_askDoesNotStop(t *testing.T) {
	out := preToolUseResults{}.Ask("confirm")
	if out.(preToolUseOutput).Stop() {
		t.Fatal("ask should not stop")
	}
}
