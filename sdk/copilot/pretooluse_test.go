package copilot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
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
	ev, err := codec.Decode([]byte(copilotPreToolUse))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(PreToolUse)
	if !ok || pre.SessionID != "s1" || pre.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if pre.NativeToolName() != "bash" || pre.ShellCommand() != "rm -rf /" {
		t.Fatalf("bad tool: name=%q shell=%q", pre.NativeToolName(), pre.ShellCommand())
	}

	out, code, err := codec.Encode(preToolResults{}.Deny("destructive command"))
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
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(PreToolUse)
	if !ok || pre.NativeToolName() != "Bash" || pre.ShellCommand() != "ls -la" {
		t.Fatalf("PreToolUse=%+v", ev)
	}
}

func TestDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{"command":"ls"}}`)
	ev, err := codec.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	got := pre.ToolInput.Raw()
	got[0] = 'X'
	if bytes.Equal(pre.ToolInput.Raw(), got) {
		t.Fatal("ToolInput.Raw() did not return a defensive copy")
	}
}

func TestEncode_PreToolAllowModifiedArgs(t *testing.T) {
	out, code, err := codec.Encode(preToolResults{}.Allow().WithModifiedArgs(map[string]any{"command": "echo safe"}))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permission_decision":"allow"`) || !strings.Contains(string(out), `"modified_args"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := codec.Encode(preToolResults{}.Noop())
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestToolInput_AsBash(t *testing.T) {
	ev, err := codec.Decode([]byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Bash","tool_input":{"command":"ls -la"}}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	input, ok := pre.Input().AsBash()
	if !ok || input.Command != "ls -la" {
		t.Fatalf("AsBash = %+v, %v", input, ok)
	}
}
