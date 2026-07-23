package pretooluse

import (
	"bytes"
	"encoding/json"
	"maps"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

const copilotPreToolUse = `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "bash",
  "tool_input": {"command": "rm -rf /"}
}`

func TestDecodeEncode_PreToolDeny(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(copilotPreToolUse))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(Event)
	if !ok || pre.SessionID != "s1" || pre.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if pre.NativeToolName() != "bash" || pre.ShellCommand() != "rm -rf /" {
		t.Fatalf("bad tool: name=%q shell=%q", pre.NativeToolName(), pre.ShellCommand())
	}

	out, code, err := results{}.Deny("destructive command").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Decision string `json:"permission_decision"`
		Reason   string `json:"permission_decision_reason"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Decision != "deny" || got.Reason != "destructive command" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_VSCodePreToolBash(t *testing.T) {
	raw := `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "Bash",
  "tool_input": {"command": "ls -la"}
}`
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(Event)
	if !ok || pre.NativeToolName() != "Bash" || pre.ShellCommand() != "ls -la" {
		t.Fatalf("PreToolUse=%+v", ev)
	}
}

func TestDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{"command":"ls"}}`)
	ev, err := runtime.Codec.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(Event)
	got := pre.ToolInput.Raw()
	got[0] = 'X'
	if bytes.Equal(pre.ToolInput.Raw(), got) {
		t.Fatal("ToolInput.Raw() did not return a defensive copy")
	}
}

func TestEncode_PreToolAllowModifiedArgs(t *testing.T) {
	out, code, err := results{}.Allow().WithModifiedArgs(map[string]any{"command": "echo safe"}).Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permission_decision":"allow"`) || !strings.Contains(string(out), `"modified_args"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out := results{}.Noop()
	if !out.IsZero() {
		t.Fatal("noop should be zero")
	}
}

func TestMerge_PreTool_denyBeatsAllowAndStops(t *testing.T) {
	a := results{}.Allow()
	b := results{}.Deny("blocked")
	merged, warnings, err := a.Merge(b.(run.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if out.decision != event.DecisionDeny || out.reason != "blocked" {
		t.Fatalf("got decision=%q reason=%q", out.decision, out.reason)
	}
	if !merged.Stop() {
		t.Fatal("deny should stop")
	}
}

func TestMerge_PreTool_modifiedArgsOverwriteWarns(t *testing.T) {
	first := map[string]any{"cmd": "a"}
	second := map[string]any{"cmd": "b"}
	orig := maps.Clone(first)
	a := results{}.Allow().WithModifiedArgs(first)
	b := results{}.Allow().WithModifiedArgs(second)
	merged, warnings, err := a.Merge(b.(run.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] != "modified_args: overwritten by later handler" {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if out.modifiedArgs["cmd"] != "b" {
		t.Fatalf("modifiedArgs = %v", out.modifiedArgs)
	}
	if !maps.Equal(first, orig) {
		t.Fatalf("caller map mutated")
	}
}

func TestToolInput_AsBash(t *testing.T) {
	ev, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Bash","tool_input":{"command":"ls -la"}}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(Event)
	input, ok := pre.Input().AsBash()
	if !ok || input.Command != "ls -la" {
		t.Fatalf("AsBash = %+v, %v", input, ok)
	}
}

func init() {
	Register(runtime.Codec)
}
