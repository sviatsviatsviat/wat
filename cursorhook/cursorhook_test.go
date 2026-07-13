package cursorhook_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cursorhook"
	"github.com/sviatsviatsviat/wat/cursorhook/tools"
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
	ev, err := cursorhook.Decode([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(cursorhook.BeforeShellExecution)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.Session() != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}
	if !bytes.Equal(cursorhook.RawBytes(ev), []byte(cursorShell)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := cursorhook.Encode(cursorhook.EventBeforeShellExecution, cursorhook.PermissionOutput{
		Decision:     cursorhook.DecisionDeny,
		AgentMessage: "force push blocked",
	})
	if err != nil {
		t.Fatal(err)
	}
	if code != cursorhook.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursorhook.PermissionDenyExit)
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
	ev, err := cursorhook.Decode([]byte(cursorStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(cursorhook.Stop)
	if !ok {
		t.Fatalf("want Stop, got %T", ev)
	}
	if stop.Status != "error" || stop.LoopCount != 1 {
		t.Fatalf("bad stop: %+v", stop)
	}

	out, code, err := cursorhook.Encode(cursorhook.EventStop, cursorhook.StopOutput{
		FollowUpMessage: "retry with fixed creds",
	})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestDecode_AfterFileEdit(t *testing.T) {
	ev, err := cursorhook.Decode([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	edit, ok := ev.(cursorhook.AfterFileEdit)
	if !ok {
		t.Fatalf("want AfterFileEdit, got %T", ev)
	}
	if edit.FilePath != "main.go" || len(edit.Edits) != 1 || edit.Edits[0].OldString != "foo" {
		t.Fatalf("bad edit: %+v", edit)
	}
	if !bytes.Equal(cursorhook.RawBytes(ev), []byte(cursorAfterFileEdit)) {
		t.Fatal("Raw not preserved")
	}
}

func TestDecode_Matrix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		eventHint string
		wantType  any
		check     func(t *testing.T, ev cursorhook.Event)
	}{
		{
			name:     "sessionStart",
			raw:      `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`,
			wantType: cursorhook.SessionStart{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.SessionStart)
				if e.Model != "gpt" || !e.IsBackgroundAgent {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "sessionEnd",
			raw:      `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`,
			wantType: cursorhook.SessionEnd{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.SessionEnd)
				if e.Reason != "complete" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeSubmitPrompt",
			raw:      `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`,
			wantType: cursorhook.BeforeSubmitPrompt{},
			check: func(t *testing.T, ev cursorhook.Event) {
				if ev.(cursorhook.BeforeSubmitPrompt).Prompt != "hello" {
					t.Fatal("bad prompt")
				}
			},
		},
		{
			name:     "preToolUse shell",
			raw:      `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			wantType: cursorhook.PreToolUse{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.PreToolUse)
				if e.ShellCommand() != "ls" {
					t.Fatalf("ShellCommand=%q", e.ShellCommand())
				}
			},
		},
		{
			name:     "postToolUse",
			raw:      `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`,
			wantType: cursorhook.PostToolUse{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.PostToolUse)
				if e.ToolOutput != "contents" || e.DurationMillis() != 100 {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "postToolUseFailure",
			raw:      `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`,
			wantType: cursorhook.PostToolUseFailure{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.PostToolUseFailure)
				if e.ErrorMessage != "timeout" || e.FailureType != "timeout" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeShellExecution",
			raw:      cursorShell,
			wantType: cursorhook.BeforeShellExecution{},
		},
		{
			name:     "afterShellExecution",
			raw:      `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`,
			wantType: cursorhook.AfterShellExecution{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.AfterShellExecution)
				if e.Command != "ls" || e.Output != "a\nb" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeMCPExecution",
			raw:      `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`,
			wantType: cursorhook.BeforeMCPExecution{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.BeforeMCPExecution)
				if e.URL != "https://mcp.example/mcp" {
					t.Fatalf("URL=%q", e.URL)
				}
			},
		},
		{
			name:     "afterMCPExecution",
			raw:      `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`,
			wantType: cursorhook.AfterMCPExecution{},
		},
		{
			name:     "beforeReadFile",
			raw:      `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`,
			wantType: cursorhook.BeforeReadFile{},
		},
		{
			name:     "afterFileEdit",
			raw:      cursorAfterFileEdit,
			wantType: cursorhook.AfterFileEdit{},
		},
		{
			name:     "subagentStart",
			raw:      `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`,
			wantType: cursorhook.SubagentStart{},
		},
		{
			name:     "subagentStop",
			raw:      `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`,
			wantType: cursorhook.SubagentStop{},
		},
		{
			name:     "stop",
			raw:      cursorStop,
			wantType: cursorhook.Stop{},
		},
		{
			name:     "preCompact",
			raw:      `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`,
			wantType: cursorhook.PreCompact{},
		},
		{
			name:     "afterAgentResponse",
			raw:      `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`,
			wantType: cursorhook.AfterAgentResponse{},
		},
		{
			name:     "afterAgentThought",
			raw:      `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`,
			wantType: cursorhook.AfterAgentThought{},
		},
		{
			name:     "beforeTabFileRead",
			raw:      `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`,
			wantType: cursorhook.BeforeTabFileRead{},
		},
		{
			name:     "afterTabFileEdit",
			raw:      `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`,
			wantType: cursorhook.AfterTabFileEdit{},
		},
		{
			name:     "workspaceOpen",
			raw:      `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`,
			wantType: cursorhook.WorkspaceOpen{},
		},
		{
			name:      "beforeShellExecution eventHint",
			raw:       `{"conversation_id":"c1","command":"ls","cwd":"/w"}`,
			eventHint: cursorhook.EventBeforeShellExecution,
			wantType:  cursorhook.BeforeShellExecution{},
			check: func(t *testing.T, ev cursorhook.Event) {
				e := ev.(cursorhook.BeforeShellExecution)
				if e.Command != "ls" || e.Cwd != "/w" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := cursorhook.Decode([]byte(tt.raw), cursorhook.WithEvent(tt.eventHint))
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
			if len(cursorhook.RawBytes(ev)) == 0 {
				t.Fatal("Raw empty")
			}
			if tt.check != nil {
				tt.check(t, ev)
			}
		})
	}
}

func reflectTypeName(v any) string {
	switch v.(type) {
	case cursorhook.SessionStart:
		return "SessionStart"
	case cursorhook.SessionEnd:
		return "SessionEnd"
	case cursorhook.BeforeSubmitPrompt:
		return "BeforeSubmitPrompt"
	case cursorhook.PreToolUse:
		return "PreToolUse"
	case cursorhook.PostToolUse:
		return "PostToolUse"
	case cursorhook.PostToolUseFailure:
		return "PostToolUseFailure"
	case cursorhook.BeforeShellExecution:
		return "BeforeShellExecution"
	case cursorhook.AfterShellExecution:
		return "AfterShellExecution"
	case cursorhook.BeforeMCPExecution:
		return "BeforeMCPExecution"
	case cursorhook.AfterMCPExecution:
		return "AfterMCPExecution"
	case cursorhook.BeforeReadFile:
		return "BeforeReadFile"
	case cursorhook.AfterFileEdit:
		return "AfterFileEdit"
	case cursorhook.SubagentStart:
		return "SubagentStart"
	case cursorhook.SubagentStop:
		return "SubagentStop"
	case cursorhook.Stop:
		return "Stop"
	case cursorhook.PreCompact:
		return "PreCompact"
	case cursorhook.AfterAgentResponse:
		return "AfterAgentResponse"
	case cursorhook.AfterAgentThought:
		return "AfterAgentThought"
	case cursorhook.BeforeTabFileRead:
		return "BeforeTabFileRead"
	case cursorhook.AfterTabFileEdit:
		return "AfterTabFileEdit"
	case cursorhook.WorkspaceOpen:
		return "WorkspaceOpen"
	default:
		return "unknown"
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := cursorhook.Encode(cursorhook.EventBeforeShellExecution, cursorhook.PermissionOutput{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	cont := false
	out, code, err := cursorhook.Encode(cursorhook.EventBeforeSubmitPrompt, cursorhook.BeforeSubmitPromptOutput{
		Continue:    &cont,
		UserMessage: "blocked",
	})
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
	out, code, err := cursorhook.Encode(cursorhook.EventSessionStart, cursorhook.SessionStartOutput{
		Env:               map[string]string{"K": "V"},
		AdditionalContext: "ctx",
	})
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
	out, code, err := cursorhook.Encode(cursorhook.EventBeforeTabFileRead, cursorhook.PermissionOutput{
		Decision:     cursorhook.DecisionDeny,
		AgentMessage: "no tab reads",
	})
	if err != nil {
		t.Fatal(err)
	}
	if code != cursorhook.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursorhook.PermissionDenyExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionUpdatedInputEmptyEventName(t *testing.T) {
	out, code, err := cursorhook.Encode("", cursorhook.PermissionOutput{
		Decision:     cursorhook.DecisionAllow,
		UpdatedInput: map[string]any{"command": "ls"},
	})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"updated_input"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestMux_Serve_BeforeShellDeny(t *testing.T) {
	mux := cursorhook.NewMux()
	cursorhook.On(mux, func(ctx context.Context, ev cursorhook.BeforeShellExecution) (cursorhook.PermissionOutput, error) {
		return cursorhook.PermissionOutput{
			Decision:     cursorhook.DecisionDeny,
			AgentMessage: "blocked",
		}, nil
	})
	var stdout bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(cursorShell), &stdout, &bytes.Buffer{})
	if code != cursorhook.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursorhook.PermissionDenyExit)
	}
	if !strings.Contains(stdout.String(), `"permission":"deny"`) {
		t.Fatalf("bad stdout: %s", stdout.String())
	}
}

func TestMux_OnDuplicatePanics(t *testing.T) {
	mux := cursorhook.NewMux()
	cursorhook.On(mux, func(ctx context.Context, ev cursorhook.BeforeShellExecution) (cursorhook.PermissionOutput, error) {
		return cursorhook.PermissionOutput{}, nil
	})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate handler registration")
		}
	}()
	cursorhook.On(mux, func(ctx context.Context, ev cursorhook.BeforeShellExecution) (cursorhook.PermissionOutput, error) {
		return cursorhook.PermissionOutput{}, nil
	})
}

func TestTools_ShellInput(t *testing.T) {
	input, err := tools.ToolInputAs[tools.ShellInput]([]byte(`{"command":"ls"}`))
	if err != nil {
		t.Fatal(err)
	}
	if input.Command != "ls" {
		t.Fatalf("Command=%q", input.Command)
	}
}

func TestParseHandler_roundTrip(t *testing.T) {
	raw := []byte(`{"command":"x.sh","matcher":"Shell","timeout":20,"loop_limit":3}`)
	h, err := cursorhook.ParseHandler(raw)
	if err != nil {
		t.Fatal(err)
	}
	if h.Command != "x.sh" || h.Matcher != "Shell" || h.TimeoutSeconds() != 20 || h.LoopLimit != 3 {
		t.Fatalf("handler=%+v", h)
	}
	out, err := cursorhook.MarshalHandler(h)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"loop_limit":3`) {
		t.Fatalf("marshal=%s", out)
	}
}
