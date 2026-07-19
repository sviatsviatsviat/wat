package cursor

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

const cursorShell = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "beforeShellExecution",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "command": "git push --force",
  "cwd": "/w",
  "sandbox": false
}`
const cursorStop = `{
  "conversation_id": "c1",
  "hook_event_name": "stop",
  "cursor_version": "1.7.2",
  "status": "error",
  "loop_count": 1
}`
const cursorAfterFileEdit = `{
  "conversation_id": "c1",
  "hook_event_name": "afterFileEdit",
  "cursor_version": "1.7.2",
  "cwd": "/w",
  "file_path": "main.go",
  "edits": [
    {"old_string": "foo", "new_string": "bar"}
  ]
}`

func TestDecodeEncode_BeforeShellDeny(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(BeforeShellExecution)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.ConversationID != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}

	out, code, err := encode(EventBeforeShellExecution, permissionResults{}.Deny("force push blocked"))
	if err != nil {
		t.Fatal(err)
	}
	if code != PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, PermissionDenyExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "force push blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecodeEncode_StopFollowUp(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(Stop)
	if !ok {
		t.Fatalf("want Stop, got %T", ev)
	}
	if stop.Status != "error" || stop.LoopCount != 1 {
		t.Fatalf("bad stop: %+v", stop)
	}

	out, code, err := encode(EventStop, stopResults{}.FollowUp("retry with fixed creds"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestDecode_AfterFileEdit(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	edit, ok := ev.(AfterFileEdit)
	if !ok {
		t.Fatalf("want AfterFileEdit, got %T", ev)
	}
	if edit.FilePath != "main.go" || len(edit.Edits) != 1 || edit.Edits[0].OldString != "foo" {
		t.Fatalf("bad edit: %+v", edit)
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	raw := `{"conversation_id":"c1","command":"ls","cwd":"/w"}`
	_, err := codec.Decode([]byte(raw))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "hook_event_name") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := codec.Decode([]byte(`{"hook_event_name":"FutureEvent","conversation_id":"c1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":123}`)
	_, err := codec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestDecode_Matrix(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		wantType any
		check    func(t *testing.T, ev run.Event)
	}{
		{
			name:     "sessionStart",
			raw:      `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`,
			wantType: SessionStart{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(SessionStart)
				if e.Model != "gpt" || !e.IsBackgroundAgent {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "sessionEnd",
			raw:      `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`,
			wantType: SessionEnd{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(SessionEnd)
				if e.Reason != "complete" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeSubmitPrompt",
			raw:      `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`,
			wantType: BeforeSubmitPrompt{},
			check: func(t *testing.T, ev run.Event) {
				if ev.(BeforeSubmitPrompt).Prompt != "hello" {
					t.Fatal("bad prompt")
				}
			},
		},
		{
			name:     "preToolUse shell",
			raw:      `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			wantType: PreToolUse{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(PreToolUse)
				if e.ShellCommand() != "ls" {
					t.Fatalf("ShellCommand=%q", e.ShellCommand())
				}
			},
		},
		{
			name:     "postToolUse",
			raw:      `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`,
			wantType: PostToolUse{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(PostToolUse)
				if e.ToolOutput != "contents" || e.DurationMillis() != 100 {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "postToolUseFailure",
			raw:      `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`,
			wantType: PostToolUseFailure{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(PostToolUseFailure)
				if e.ErrorMessage != "timeout" || e.FailureType != "timeout" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeShellExecution",
			raw:      cursorShell,
			wantType: BeforeShellExecution{},
		},
		{
			name:     "afterShellExecution",
			raw:      `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`,
			wantType: AfterShellExecution{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(AfterShellExecution)
				if e.Command != "ls" || e.Output != "a\nb" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeMCPExecution",
			raw:      `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`,
			wantType: BeforeMCPExecution{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(BeforeMCPExecution)
				if e.URL != "https://mcp.example/mcp" {
					t.Fatalf("URL=%q", e.URL)
				}
			},
		},
		{
			name:     "afterMCPExecution",
			raw:      `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`,
			wantType: AfterMCPExecution{},
		},
		{
			name:     "beforeReadFile",
			raw:      `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`,
			wantType: BeforeReadFile{},
		},
		{
			name:     "afterFileEdit",
			raw:      cursorAfterFileEdit,
			wantType: AfterFileEdit{},
		},
		{
			name:     "subagentStart",
			raw:      `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`,
			wantType: SubagentStart{},
		},
		{
			name:     "subagentStop",
			raw:      `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`,
			wantType: SubagentStop{},
		},
		{
			name:     "stop",
			raw:      cursorStop,
			wantType: Stop{},
		},
		{
			name:     "preCompact",
			raw:      `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`,
			wantType: PreCompact{},
		},
		{
			name:     "afterAgentResponse",
			raw:      `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`,
			wantType: AfterAgentResponse{},
		},
		{
			name:     "afterAgentThought",
			raw:      `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`,
			wantType: AfterAgentThought{},
		},
		{
			name:     "beforeTabFileRead",
			raw:      `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`,
			wantType: BeforeTabFileRead{},
		},
		{
			name:     "afterTabFileEdit",
			raw:      `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`,
			wantType: AfterTabFileEdit{},
		},
		{
			name:     "workspaceOpen",
			raw:      `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`,
			wantType: WorkspaceOpen{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := codec.Decode([]byte(tt.raw))
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantType != nil {
				want := reflectTypeName(tt.wantType)
				got := reflectTypeName(ev)
				if want != got {
					t.Fatalf("type %s != %s", got, want)
				}
			}
			if ev.EventName() == "" {
				t.Fatal("EventName empty")
			}
			if tt.check != nil {
				tt.check(t, ev)
			}
		})
	}
}

func reflectTypeName(v any) string {
	switch v.(type) {
	case SessionStart:
		return "SessionStart"
	case SessionEnd:
		return "SessionEnd"
	case BeforeSubmitPrompt:
		return "BeforeSubmitPrompt"
	case PreToolUse:
		return "PreToolUse"
	case PostToolUse:
		return "PostToolUse"
	case PostToolUseFailure:
		return "PostToolUseFailure"
	case BeforeShellExecution:
		return "BeforeShellExecution"
	case AfterShellExecution:
		return "AfterShellExecution"
	case BeforeMCPExecution:
		return "BeforeMCPExecution"
	case AfterMCPExecution:
		return "AfterMCPExecution"
	case BeforeReadFile:
		return "BeforeReadFile"
	case AfterFileEdit:
		return "AfterFileEdit"
	case SubagentStart:
		return "SubagentStart"
	case SubagentStop:
		return "SubagentStop"
	case Stop:
		return "Stop"
	case PreCompact:
		return "PreCompact"
	case AfterAgentResponse:
		return "AfterAgentResponse"
	case AfterAgentThought:
		return "AfterAgentThought"
	case BeforeTabFileRead:
		return "BeforeTabFileRead"
	case AfterTabFileEdit:
		return "AfterTabFileEdit"
	case WorkspaceOpen:
		return "WorkspaceOpen"
	default:
		return "unknown"
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := encode(EventBeforeShellExecution, nil)
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	out, code, err := encode(EventBeforeSubmitPrompt, beforeSubmitPromptResults{}.Block("blocked"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["continue"] != false || got["user_message"] != "blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_SessionStartEnv(t *testing.T) {
	out, code, err := encode(EventSessionStart, sessionStartResults{}.Noop().
		WithEnv(map[string]string{"K": "V"}).
		WithAdditionalContext("ctx"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Env map[string]string `json:"env"`
		Ctx string            `json:"additional_context"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Env["K"] != "V" || got.Ctx != "ctx" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_TabFileReadDeny(t *testing.T) {
	out, code, err := encode(EventBeforeTabFileRead, permissionResults{}.Deny("no tab reads"))
	if err != nil {
		t.Fatal(err)
	}
	if code != PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, PermissionDenyExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionUpdatedInputEmptyEventName(t *testing.T) {
	out, code, err := encode("", permissionResults{}.Allow().WithUpdatedInput(map[string]any{"command": "ls"}))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"updated_input"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestToolInput_AsShell(t *testing.T) {
	ev, err := codec.Decode([]byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	input, ok := pre.ToolInput.AsShell()
	if !ok || input.Command != "ls" {
		t.Fatalf("AsShell = %+v, %v", input, ok)
	}
}
