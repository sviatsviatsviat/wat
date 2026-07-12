package agenthooks

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Design §3.2 ⚠ verification checklist (re-verify against live Claude docs):
//
// Decoded into unified sub-structs:
//   - SessionEnd.reason → Life.Reason
//   - PostToolUseFailure.error → Result.Error
//   - SubagentStop.stop_hook_active → Turn.StopHookActive
//   - Stop.stop_hook_active → Turn.StopHookActive
//   - StopFailure.error_type → Note.Type
//   - Notification.notification_type → Note.Type
//   - PreCompact.custom_instructions → Compact.CustomInstructions
//
// Preserved only via Raw (KindOther):
//   - PostToolBatch, TaskCreated.task, TaskCompleted.task
//   - ConfigChange.source, CwdChanged.new_cwd/old_cwd, FileChanged.file_path
//   - WorktreeRemove.worktree_path, PostCompact.trigger
//   - Elicitation.*, ElicitationResult.*
//   - SubagentStop.agent_transcript_path
//
// Encode ⚠ mapped via unified Result:
//   - PreToolUseOutput.updatedInput → Result.UpdatedInput
//   - PostToolUseOutput.updatedToolOutput → Result.UpdatedOutput
//   - PermissionRequestOutput.message → Result.Reason; Interrupt deferred

func TestCodecFor(t *testing.T) {
	tests := []struct {
		name    string
		dialect Dialect
		wantErr bool
	}{
		{name: "claude", dialect: Claude},
		{name: "copilot", dialect: Copilot},
		{name: "cursor", dialect: Cursor},
		{name: "unknown", dialect: Unknown, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec, err := CodecFor(tt.dialect)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if codec.Dialect() != tt.dialect {
				t.Fatalf("Dialect() = %v, want %v", codec.Dialect(), tt.dialect)
			}
		})
	}
}

func TestClaudeDecodeEncode_PreToolDeny(t *testing.T) {
	c := &ClaudeCodec{Getenv: func(string) string { return "" }}
	ev, err := c.Decode([]byte(claudePreToolUse), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != KindPreTool || ev.Session != "abc123" || ev.Cwd != "/home/user/proj" {
		t.Fatalf("bad event: %+v", ev)
	}
	if ev.Tool == nil || ev.Tool.Name != ToolBash || ev.Tool.Native != "Bash" || ev.Tool.Shell != "rm -rf /tmp/build" {
		t.Fatalf("bad tool: %+v", ev.Tool)
	}
	if !bytes.Equal(ev.Raw, []byte(claudePreToolUse)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := c.Encode(ev, Deny("destructive command"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		HSO struct {
			Event    string `json:"hookEventName"`
			Decision string `json:"permissionDecision"`
			Reason   string `json:"permissionDecisionReason"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.HSO.Event != "PreToolUse" || got.HSO.Decision != "deny" || got.HSO.Reason != "destructive command" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeDecodeEncode_PreToolUpdatedInput(t *testing.T) {
	c := &ClaudeCodec{Getenv: func(string) string { return "" }}
	ev, err := c.Decode([]byte(claudePreToolUse), "")
	if err != nil {
		t.Fatal(err)
	}

	out, code, err := c.Encode(ev, Result{UpdatedInput: map[string]any{"command": "ls -la"}})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		HSO struct {
			Decision string         `json:"permissionDecision"`
			Input    map[string]any `json:"updatedInput"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.HSO.Decision != "allow" {
		t.Fatalf("permissionDecision = %q, want allow", got.HSO.Decision)
	}
	if got.HSO.Input["command"] != "ls -la" {
		t.Fatalf("updatedInput = %+v", got.HSO.Input)
	}
}

func TestClaudeEncode_StopFollowUp(t *testing.T) {
	c := &ClaudeCodec{}
	stopEv := &Event{Agent: Claude, Kind: KindStop, Name: "Stop"}
	out, code, err := c.Encode(stopEv, Result{FollowUp: "run the tests"})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestClaudeEncode_SubagentStopFollowUp(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindSubagentStop, Name: "SubagentStop"}
	out, code, err := c.Encode(ev, Result{FollowUp: "finish the review"})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "finish the review") {
		t.Fatalf("bad subagent stop output: %s", out)
	}
}

func TestClaudeEncode_SessionStartEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.sh")
	var written []byte
	c := &ClaudeCodec{
		Getenv: func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return envPath
			}
			return ""
		},
		AppendFile: func(path string, data []byte) error {
			if path != envPath {
				t.Fatalf("unexpected path %q", path)
			}
			written = append(written, data...)
			return os.WriteFile(path, written, 0o644)
		},
	}
	ev := &Event{Agent: Claude, Kind: KindSessionStart, Name: "SessionStart"}
	out, code, err := c.Encode(ev, Result{Env: map[string]string{"FOO": "bar", "BAZ": "qux"}})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if out != nil {
		t.Fatalf("env-only result should produce no stdout, got %q", out)
	}
	wantLines := []string{`export FOO='bar'`, `export BAZ='qux'`}
	got := string(written)
	for _, line := range wantLines {
		if !strings.Contains(got, line) {
			t.Fatalf("env file = %q, want lines %v", got, wantLines)
		}
	}
	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range wantLines {
		if !strings.Contains(string(data), line) {
			t.Fatalf("file contents = %q, want lines %v", data, wantLines)
		}
	}

	t.Run("shell metacharacters preserved", func(t *testing.T) {
		written = nil
		ev := &Event{Agent: Claude, Kind: KindSessionStart, Name: "SessionStart"}
		unsafe := "$(rm -rf /) `whoami` $HOME"
		_, code, err := c.Encode(ev, Result{Env: map[string]string{"UNSAFE": unsafe}})
		if err != nil || code != 0 {
			t.Fatalf("encode: %v code=%d", err, code)
		}
		want := "export UNSAFE='$(rm -rf /) `whoami` $HOME'\n"
		if string(written) != want {
			t.Fatalf("env file = %q, want %q", written, want)
		}
	})
}

func TestClaudeEncode_ZeroResult(t *testing.T) {
	c := &ClaudeCodec{}
	ev, _ := c.Decode([]byte(claudePreToolUse), "")
	out, code, err := c.Encode(ev, Result{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestClaudeDecode_UnknownEvent(t *testing.T) {
	raw := []byte(`{"session_id":"s1","hook_event_name":"Setup","cwd":"/w"}`)
	c := &ClaudeCodec{}
	ev, err := c.Decode(raw, "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != KindOther || ev.Name != "Setup" {
		t.Fatalf("Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if !bytes.Equal(ev.Raw, raw) {
		t.Fatal("Raw not preserved")
	}
}

func TestClaudeDecode_Matrix(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		kind  Kind
		check func(t *testing.T, ev *Event)
	}{
		{
			name: "SessionStart",
			raw:  `{"session_id":"s","hook_event_name":"SessionStart","source":"startup","model":"claude-3"}`,
			kind: KindSessionStart,
			check: func(t *testing.T, ev *Event) {
				if ev.Life == nil || ev.Life.Source != "startup" || ev.Life.Model != "claude-3" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "SessionEnd",
			raw:  `{"session_id":"s","hook_event_name":"SessionEnd","reason":"clear"}`,
			kind: KindSessionEnd,
			check: func(t *testing.T, ev *Event) {
				if ev.Life == nil || ev.Life.Reason != "clear" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "UserPromptSubmit",
			raw:  `{"session_id":"s","hook_event_name":"UserPromptSubmit","prompt":"hello"}`,
			kind: KindUserPrompt,
			check: func(t *testing.T, ev *Event) {
				if ev.Prompt != "hello" {
					t.Fatalf("Prompt=%q", ev.Prompt)
				}
			},
		},
		{
			name: "PreToolUse",
			raw:  `{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			kind: KindPreTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Shell != "ls" || ev.Tool.Name != ToolBash {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "PostToolUse",
			raw:  `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`,
			kind: KindPostTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Result == nil || ev.Result.Text != "file contents" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "PostToolUseFailure",
			raw:  `{"session_id":"s","hook_event_name":"PostToolUseFailure","tool_name":"Bash","error":"timeout"}`,
			kind: KindPostToolFailure,
			check: func(t *testing.T, ev *Event) {
				if ev.Result == nil || ev.Result.Error != "timeout" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "PermissionRequest",
			raw:  `{"session_id":"s","hook_event_name":"PermissionRequest","tool_name":"Write","tool_use_id":"t2"}`,
			kind: KindPermissionRequest,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Name != ToolWrite {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "SubagentStart",
			raw:  `{"session_id":"s","hook_event_name":"SubagentStart","agent_id":"a1","agent_type":"reviewer"}`,
			kind: KindSubagentStart,
			check: func(t *testing.T, ev *Event) {
				if ev.Subagent == nil || ev.Subagent.ID != "a1" || ev.Subagent.Type != "reviewer" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "SubagentStop",
			raw:  `{"session_id":"s","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"worker","stop_hook_active":true,"last_assistant_message":"done"}`,
			kind: KindSubagentStop,
			check: func(t *testing.T, ev *Event) {
				if ev.Subagent == nil || ev.Subagent.ID != "a1" || ev.Subagent.Type != "worker" || ev.Subagent.Summary != "done" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
				if ev.Turn == nil || !ev.Turn.StopHookActive || ev.Turn.LastAssistantMessage != "done" {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "Stop",
			raw:  `{"session_id":"s","hook_event_name":"Stop","stop_hook_active":false,"last_assistant_message":"bye"}`,
			kind: KindStop,
			check: func(t *testing.T, ev *Event) {
				if ev.Turn == nil || ev.Turn.LastAssistantMessage != "bye" {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "PreCompact",
			raw:  `{"session_id":"s","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":"keep tests"}`,
			kind: KindPreCompact,
			check: func(t *testing.T, ev *Event) {
				if ev.Compact == nil || ev.Compact.Trigger != "manual" || ev.Compact.CustomInstructions != "keep tests" {
					t.Fatalf("Compact=%+v", ev.Compact)
				}
			},
		},
		{
			name: "Notification",
			raw:  `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`,
			kind: KindNotification,
			check: func(t *testing.T, ev *Event) {
				if ev.Note == nil || ev.Note.Type != "idle_prompt" || ev.Note.Message != "waiting" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "StopFailure",
			raw:  `{"session_id":"s","hook_event_name":"StopFailure","error_type":"rate_limit","message":"slow down"}`,
			kind: KindAgentError,
			check: func(t *testing.T, ev *Event) {
				if ev.Note == nil || ev.Note.Type != "rate_limit" || ev.Note.Message != "slow down" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
	}
	c := &ClaudeCodec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := c.Decode([]byte(tt.raw), "")
			if err != nil {
				t.Fatal(err)
			}
			if ev.Kind != tt.kind {
				t.Fatalf("Kind=%v, want %v", ev.Kind, tt.kind)
			}
			if ev.Name != tt.name {
				t.Fatalf("Name=%q", ev.Name)
			}
			tt.check(t, ev)
		})
	}
}

func TestClaudeEncode_SessionStartEnvInvalidKey(t *testing.T) {
	c := &ClaudeCodec{
		Getenv: func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return "/tmp/env.sh"
			}
			return ""
		},
	}
	ev := &Event{Agent: Claude, Kind: KindSessionStart, Name: "SessionStart"}
	_, _, err := c.Encode(ev, Result{Env: map[string]string{"FOO\nBAR": "value"}})
	if err == nil {
		t.Fatal("expected error for invalid env key")
	}
	if !strings.Contains(err.Error(), "invalid env key") {
		t.Fatalf("error = %v", err)
	}
}

func TestClaudeEncode_SessionStartEnvInvalidKeyUnsetFile(t *testing.T) {
	c := &ClaudeCodec{Getenv: func(string) string { return "" }}
	ev := &Event{Agent: Claude, Kind: KindSessionStart, Name: "SessionStart"}
	_, code, err := c.Encode(ev, Result{Env: map[string]string{"FOO\nBAR": "value"}})
	if err != nil || code != 0 {
		t.Fatalf("unset CLAUDE_ENV_FILE should no-op invalid keys, got err=%v code=%d", err, code)
	}
}

func TestClaudeEncode_NilEvent(t *testing.T) {
	c := &ClaudeCodec{}
	_, _, err := c.Encode(nil, Deny("nope"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "nil event") {
		t.Fatalf("error = %v", err)
	}
}

func TestClaudeDecode_LongTailKindOther(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{
			name: "TaskCreated",
			raw:  `{"session_id":"s","hook_event_name":"TaskCreated","task":{"id":"t1"}}`,
		},
		{
			name: "CwdChanged",
			raw:  `{"session_id":"s","hook_event_name":"CwdChanged","new_cwd":"/new","old_cwd":"/old"}`,
		},
		{
			name: "Elicitation",
			raw:  `{"session_id":"s","hook_event_name":"Elicitation","server_name":"srv","message":"confirm?"}`,
		},
	}
	c := &ClaudeCodec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.raw)
			ev, err := c.Decode(raw, "")
			if err != nil {
				t.Fatal(err)
			}
			if ev.Kind != KindOther || ev.Name != tt.name {
				t.Fatalf("Kind=%v Name=%q", ev.Kind, ev.Name)
			}
			if !bytes.Equal(ev.Raw, raw) {
				t.Fatal("Raw not preserved")
			}
		})
	}
}

func TestClaudeDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"}}`)
	c := &ClaudeCodec{}
	ev, err := c.Decode(raw, "")
	if err != nil {
		t.Fatal(err)
	}
	rawCopy := bytes.Clone(ev.Raw)
	ev.Tool.Input[0] = '{'
	if !bytes.Equal(ev.Raw, rawCopy) {
		t.Fatal("mutating Tool.Input affected Event.Raw")
	}
}

func TestClaudeEncode_PermissionRequestAllow(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindPermissionRequest, Name: "PermissionRequest"}
	out, code, err := c.Encode(ev, Allow())
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"behavior":"allow"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeEncode_PermissionRequestAsk(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindPermissionRequest, Name: "PermissionRequest"}
	out, code, err := c.Encode(ev, Ask("needs user confirmation"))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"behavior":"ask"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestValidEnvKey(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{key: "FOO", want: true},
		{key: "_PRIVATE", want: true},
		{key: "VAR2", want: true},
		{key: "", want: false},
		{key: "2BAD", want: false},
		{key: "FOO=BAR", want: false},
		{key: "FOO\nBAR", want: false},
		{key: "FOO;BAR", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := validEnvKey(tt.key); got != tt.want {
				t.Fatalf("validEnvKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestClaudeEncode_PostToolUpdatedOutput(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindPostTool, Name: "PostToolUse"}
	text := "rewritten"
	out, code, err := c.Encode(ev, Result{UpdatedOutput: &text})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"updatedToolOutput":"rewritten"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeEncode_HaltSession(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindStop, Name: "Stop"}
	out, code, err := c.Encode(ev, Result{HaltSession: true, Reason: "policy violation"})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"continue":false`) || !strings.Contains(string(out), "policy violation") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeEncode_UserPromptBlock(t *testing.T) {
	c := &ClaudeCodec{}
	ev := &Event{Agent: Claude, Kind: KindUserPrompt, Name: "UserPromptSubmit"}
	out, code, err := c.Encode(ev, Result{BlockPrompt: true, Reason: "blocked prompt"})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeDecode_InvalidJSON(t *testing.T) {
	c := &ClaudeCodec{}
	_, err := c.Decode([]byte("not json"), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "claude: decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestNewToolCall_BashShell(t *testing.T) {
	tc := newToolCall("Bash", []byte(`{"command":"echo hi","description":"test"}`), "id1")
	if tc.Name != ToolBash || tc.Native != "Bash" || tc.Shell != "echo hi" || tc.ID != "id1" {
		t.Fatalf("ToolCall=%+v", tc)
	}
}

func TestRawToText(t *testing.T) {
	if got := rawToText([]byte(`"hello"`)); got != "hello" {
		t.Fatalf("string = %q", got)
	}
	if got := rawToText([]byte(`{"k":"v"}`)); got != `{"k":"v"}` {
		t.Fatalf("object = %q", got)
	}
}
