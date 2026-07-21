package claude

import (
	"strings"
	"testing"
)

func TestEncode_PreToolDeny(t *testing.T) {
	out, code, err := codec.Encode(preToolUseResults{}.Deny("destructive command"))
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
	out, code, err := codec.Encode(preToolUseResults{}.Noop())
	if err != nil || out != nil || code != SuccessExit {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
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
