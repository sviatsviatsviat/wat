package claude_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
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

func TestDecode_PreToolUse(t *testing.T) {
	ev, err := claude.DecodeForTest([]byte(preToolUsePayload))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(claude.PreToolUse)
	if !ok || pre.ToolName != "Bash" || pre.SessionID != "abc123" {
		t.Fatalf("bad event: %+v", ev)
	}
}

func TestDecode_Matrix(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{name: "SessionStart", raw: `{"session_id":"s","hook_event_name":"SessionStart","source":"startup","model":"claude-3"}`},
		{name: "SessionEnd", raw: `{"session_id":"s","hook_event_name":"SessionEnd","reason":"clear"}`},
		{name: "UserPromptSubmit", raw: `{"session_id":"s","hook_event_name":"UserPromptSubmit","prompt":"hello"}`},
		{name: "PreToolUse", raw: `{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"},"tool_use_id":"t1"}`},
		{name: "PostToolUse", raw: `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`},
		{name: "PostToolUseFailure", raw: `{"session_id":"s","hook_event_name":"PostToolUseFailure","tool_name":"Bash","error":"timeout"}`},
		{name: "PermissionRequest", raw: `{"session_id":"s","hook_event_name":"PermissionRequest","tool_name":"Write","tool_use_id":"t2"}`},
		{name: "SubagentStart", raw: `{"session_id":"s","hook_event_name":"SubagentStart","agent_id":"a1","agent_type":"reviewer"}`},
		{name: "SubagentStop", raw: `{"session_id":"s","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"worker","stop_hook_active":true,"last_assistant_message":"done"}`},
		{name: "Stop", raw: `{"session_id":"s","hook_event_name":"Stop","stop_hook_active":false,"last_assistant_message":"bye"}`},
		{name: "PreCompact", raw: `{"session_id":"s","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":"keep tests"}`},
		{name: "Notification", raw: `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`},
		{name: "StopFailure", raw: `{"session_id":"s","hook_event_name":"StopFailure","error_type":"rate_limit","message":"slow down"}`},
		{name: "Setup", raw: `{"session_id":"s","hook_event_name":"Setup","trigger":"init"}`},
		{name: "UserPromptExpansion", raw: `{"session_id":"s","hook_event_name":"UserPromptExpansion","expansion_type":"slash_command","command_name":"foo"}`},
		{name: "PostToolBatch", raw: `{"session_id":"s","hook_event_name":"PostToolBatch"}`},
		{name: "PermissionDenied", raw: `{"session_id":"s","hook_event_name":"PermissionDenied","tool_name":"Bash"}`},
		{name: "TaskCreated", raw: `{"session_id":"s","hook_event_name":"TaskCreated","task":{"id":"t1"}}`},
		{name: "TaskCompleted", raw: `{"session_id":"s","hook_event_name":"TaskCompleted","task":{"id":"t1"}}`},
		{name: "TeammateIdle", raw: `{"session_id":"s","hook_event_name":"TeammateIdle"}`},
		{name: "MessageDisplay", raw: `{"session_id":"s","hook_event_name":"MessageDisplay","turn_id":"1","message_id":"m1","index":0,"final":true,"delta":"hi"}`},
		{name: "InstructionsLoaded", raw: `{"session_id":"s","hook_event_name":"InstructionsLoaded","file_path":"/f","memory_type":"Project","load_reason":"session_start"}`},
		{name: "ConfigChange", raw: `{"session_id":"s","hook_event_name":"ConfigChange","source":"user_settings"}`},
		{name: "CwdChanged", raw: `{"session_id":"s","hook_event_name":"CwdChanged","new_cwd":"/new","old_cwd":"/old"}`},
		{name: "FileChanged", raw: `{"session_id":"s","hook_event_name":"FileChanged","file_path":"/f.go"}`},
		{name: "WorktreeCreate", raw: `{"session_id":"s","hook_event_name":"WorktreeCreate"}`},
		{name: "WorktreeRemove", raw: `{"session_id":"s","hook_event_name":"WorktreeRemove","worktree_path":"/wt"}`},
		{name: "PostCompact", raw: `{"session_id":"s","hook_event_name":"PostCompact","trigger":"auto"}`},
		{name: "Elicitation", raw: `{"session_id":"s","hook_event_name":"Elicitation","server_name":"srv","message":"confirm?"}`},
		{name: "ElicitationResult", raw: `{"session_id":"s","hook_event_name":"ElicitationResult","server_name":"srv","action":"accept"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := claude.DecodeForTest([]byte(tt.raw))
			if err != nil {
				t.Fatal(err)
			}
			if ev.EventName() != tt.name {
				t.Fatalf("EventName() = %q, want %q", ev.EventName(), tt.name)
			}
		})
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := claude.DecodeForTest([]byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := claude.DecodeForTest([]byte(`{"session_id":"s1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, claude.ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	t.Run("envelope", func(t *testing.T) {
		_, err := claude.DecodeForTest([]byte("not json"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claude.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("typed event", func(t *testing.T) {
		raw := []byte(`{"session_id":"s1","hook_event_name":"PreToolUse","cwd":"/w","tool_name":123}`)
		_, err := claude.DecodeForTest(raw)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claude.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestEncode_PreToolDeny(t *testing.T) {
	out, err := claude.Encode("PreToolUse", claude.PreToolUseResultsForTest().Deny("destructive command"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_UserPromptBlock(t *testing.T) {
	out, err := claude.Encode("UserPromptSubmit", claude.UserPromptSubmitResultsForTest().Block("blocked prompt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopBlock(t *testing.T) {
	out, err := claude.Encode("Stop", claude.StopResultsForTest().FollowUp("run the tests"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_SessionStartEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.sh")
	var written []byte
	out, err := claude.Encode("SessionStart", claude.SessionStartResultsForTest().Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}),
		claude.WithGetenv(func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return envPath
			}
			return ""
		}),
		claude.WithAppendFile(func(path string, data []byte) error {
			written = append(written, data...)
			return os.WriteFile(path, written, 0o644)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("env-only result should produce no stdout, got %q", out)
	}
	got := string(written)
	for _, line := range []string{`export FOO='bar'`, `export BAZ='qux'`} {
		if !strings.Contains(got, line) {
			t.Fatalf("env file = %q, want lines containing %s", got, line)
		}
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, err := claude.Encode("PreToolUse", nil)
	if err != nil || out != nil {
		t.Fatalf("zero output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, err := claude.Encode(claude.EventPostToolUse, claude.PreToolUseResultsForTest().Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestEncode_PointerOutput(t *testing.T) {
	deny := claude.PreToolUseResultsForTest().Deny("blocked")
	out, err := claude.Encode("PreToolUse", &deny)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}

	var typedNil *claude.PreToolUseOutput
	out, err = claude.Encode("PreToolUse", typedNil)
	if err != nil || out != nil {
		t.Fatalf("nil pointer output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, err := claude.Encode("PermissionRequest", claude.PermissionRequestResultsForTest().Deny("policy").WithInterrupt(true))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), `"continue":false`) {
		t.Fatalf("interrupt must not set top-level continue: %s", out)
	}
	if !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestWriteEnvFile_InvalidKey(t *testing.T) {
	err := claude.WriteEnvFile(
		map[string]string{"FOO\nBAR": "value"},
		func(string) string { return "/tmp/env.sh" },
		nil,
	)
	if err == nil {
		t.Fatal("expected error for invalid env key")
	}
}

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

func TestServe_FailPolicy(t *testing.T) {
	run.Reset()
	t.Cleanup(run.Reset)
	claude.OnPreToolUse(func(ctx context.Context, hook run.Hook[claude.PreToolUse], _ claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return nil, context.Canceled
	})

	code := run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &bytes.Buffer{}, &bytes.Buffer{}, claude.WithFailPolicy(claude.FailOpen))
	if code != claude.HandlerErrorExit {
		t.Fatalf("FailOpen exit = %d, want %d", code, claude.HandlerErrorExit)
	}
	run.Reset()
	claude.OnPreToolUse(func(ctx context.Context, hook run.Hook[claude.PreToolUse], _ claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return nil, context.Canceled
	})
	code = run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &bytes.Buffer{}, &bytes.Buffer{}, claude.WithFailPolicy(claude.FailBlock))
	if code != claude.FailBlockExit {
		t.Fatalf("FailBlock exit = %d, want %d", code, claude.FailBlockExit)
	}
}
