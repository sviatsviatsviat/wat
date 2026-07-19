package claude_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

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

func TestServe_PreToolDeny(t *testing.T) {
	run.Reset()
	t.Cleanup(run.Reset)
	claude.OnPreToolUse(func(ctx context.Context, hook run.Hook[claude.PreToolUse], r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return r.Deny("destructive command"), nil
	})

	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permissionDecision":"deny"`) {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_HandlerErrorExit(t *testing.T) {
	run.Reset()
	t.Cleanup(run.Reset)
	claude.OnPreToolUse(func(ctx context.Context, hook run.Hook[claude.PreToolUse], _ claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return nil, context.Canceled
	})

	code := run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &bytes.Buffer{}, &bytes.Buffer{})
	if code != claude.HandlerErrorExit {
		t.Fatalf("exit = %d, want %d", code, claude.HandlerErrorExit)
	}
}
