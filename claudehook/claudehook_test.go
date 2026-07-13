package claudehook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/claudehook"
	"github.com/sviatsviatsviat/wat/claudehook/tools"
)

const claudePreToolUse = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "PreToolUse",
  "tool_name": "Bash",
  "tool_use_id": "tu_1",
  "tool_input": {"command": "rm -rf /tmp/build", "description": "clean"}
}`

func TestDecodeEncode_PreToolDeny(t *testing.T) {
	ev, err := claudehook.Decode([]byte(claudePreToolUse))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(claudehook.PreToolUse)
	if !ok || pre.ToolName != "Bash" || pre.SessionID != "abc123" {
		t.Fatalf("bad event: %+v", ev)
	}

	out, err := claudehook.Encode("PreToolUse", claudehook.PreToolUseOutput{
		Decision: claudehook.DecisionDeny,
		Reason:   "destructive command",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestMux_Serve_PreToolDeny(t *testing.T) {
	mux := claudehook.NewMux()
	claudehook.On(mux, func(ctx context.Context, ev claudehook.PreToolUse) (claudehook.PreToolUseOutput, error) {
		return claudehook.PreToolUseOutput{
			Decision: claudehook.DecisionDeny,
			Reason:   "destructive command",
		}, nil
	})

	var stdout bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permissionDecision":"deny"`) {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestEncode_UserPromptBlock(t *testing.T) {
	out, err := claudehook.Encode("UserPromptSubmit", claudehook.UserPromptSubmitOutput{
		Block:  true,
		Reason: "blocked prompt",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopBlock(t *testing.T) {
	out, err := claudehook.Encode("Stop", claudehook.StopOutput{
		Block:  true,
		Reason: "run the tests",
	})
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
	out, err := claudehook.Encode("SessionStart", claudehook.SessionStartOutput{
		Env: map[string]string{"FOO": "bar", "BAZ": "qux"},
	},
		claudehook.WithGetenv(func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return envPath
			}
			return ""
		}),
		claudehook.WithAppendFile(func(path string, data []byte) error {
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
			ev, err := claudehook.Decode([]byte(tt.raw))
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
	raw := []byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`)
	ev, err := claudehook.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	re, ok := ev.(claudehook.RawEvent)
	if !ok {
		t.Fatalf("want RawEvent, got %T", ev)
	}
	if !bytes.Equal(re.Raw, raw) {
		t.Fatal("Raw not preserved")
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	t.Run("envelope", func(t *testing.T) {
		_, err := claudehook.Decode([]byte("not json"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claudehook.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("typed event", func(t *testing.T) {
		raw := []byte(`{"session_id":"s1","hook_event_name":"PreToolUse","cwd":"/w","tool_name":123}`)
		_, err := claudehook.Decode(raw)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claudehook.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestMux_FailPolicy(t *testing.T) {
	mux := claudehook.NewMux()
	claudehook.On(mux, func(ctx context.Context, ev claudehook.PreToolUse) (claudehook.PreToolUseOutput, error) {
		return claudehook.PreToolUseOutput{}, context.Canceled
	})

	code := mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &bytes.Buffer{}, &bytes.Buffer{}, claudehook.WithFailPolicy(claudehook.FailOpen))
	if code != 1 {
		t.Fatalf("FailOpen exit = %d, want 1", code)
	}
	code = mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &bytes.Buffer{}, &bytes.Buffer{}, claudehook.WithFailPolicy(claudehook.FailBlock))
	if code != 2 {
		t.Fatalf("FailBlock exit = %d, want 2", code)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, err := claudehook.Encode("PreToolUse", claudehook.PreToolUseOutput{})
	if err != nil || out != nil {
		t.Fatalf("zero output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, err := claudehook.Encode(claudehook.EventPostToolUse, claudehook.PreToolUseOutput{
		Decision: claudehook.DecisionAllow,
	})
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestEncode_PointerOutput(t *testing.T) {
	deny := claudehook.DecisionDeny
	out, err := claudehook.Encode("PreToolUse", &claudehook.PreToolUseOutput{
		Decision: deny,
		Reason:   "blocked",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}

	out, err = claudehook.Encode("PreToolUse", (*claudehook.PreToolUseOutput)(nil))
	if err != nil || out != nil {
		t.Fatalf("nil pointer output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, err := claudehook.Encode("PermissionRequest", claudehook.PermissionRequestOutput{
		Behavior:  "deny",
		Message:   "policy",
		Interrupt: true,
	})
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

func TestRawBytes_PreservesTypedDecode(t *testing.T) {
	raw := []byte(claudePreToolUse)
	ev, err := claudehook.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	got := claudehook.RawBytes(ev)
	if !bytes.Equal(got, raw) {
		t.Fatalf("RawBytes = %s, want original payload", got)
	}
}

func TestToolInputAs_GrepCaseInsensitive(t *testing.T) {
	input, err := tools.ToolInputAs[tools.GrepInput](json.RawMessage(`{"pattern":"foo","-i":true}`))
	if err != nil {
		t.Fatal(err)
	}
	if !input.CaseInsensitive {
		t.Fatal("expected CaseInsensitive from -i key")
	}
}

func TestToolInputAs_AskUserQuestionMultiSelect(t *testing.T) {
	input, err := tools.ToolInputAs[tools.AskUserQuestionInput](json.RawMessage(`{
		"questions":[{"question":"Pick","header":"H","options":[{"label":"A"}],"multiSelect":true}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(input.Questions) != 1 || !input.Questions[0].MultiSelect {
		t.Fatalf("MultiSelect not decoded: %+v", input.Questions[0])
	}
}

func TestToolInputAs_Bash(t *testing.T) {
	input, err := tools.ToolInputAs[tools.BashInput](json.RawMessage(`{"command":"ls -la","description":"list"}`))
	if err != nil {
		t.Fatal(err)
	}
	if input.Command != "ls -la" {
		t.Fatalf("Command = %q", input.Command)
	}
}

func TestIsMCPTool(t *testing.T) {
	server, tool, ok := tools.IsMCPTool("mcp__github__create_issue")
	if !ok || server != "github" || tool != "create_issue" {
		t.Fatalf("IsMCPTool = %q, %q, %v", server, tool, ok)
	}
}

func TestParseHandler_RoundTrip(t *testing.T) {
	raw, err := claudehook.MarshalHandler(claudehook.Handler{
		Type:    "command",
		Command: "echo hi",
		Timeout: 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := claudehook.ParseHandler(raw)
	if err != nil {
		t.Fatal(err)
	}
	if h.Type != "command" || h.Command != "echo hi" || h.Timeout != 30 {
		t.Fatalf("handler = %+v", h)
	}
	if h.TimeoutSeconds() != 30 {
		t.Fatalf("TimeoutSeconds = %d", h.TimeoutSeconds())
	}
}

func TestHandlers_EncodesMultiple(t *testing.T) {
	blobs, err := claudehook.Handlers(
		claudehook.Handler{Type: "command", Command: "a"},
		claudehook.Handler{Type: "command", Command: "b"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(blobs) != 2 {
		t.Fatalf("len = %d", len(blobs))
	}
}
