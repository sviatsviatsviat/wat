package cursor_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
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
	ev, err := cursor.DecodeForTest([]byte(cursorShell))
	if err != nil {
		t.Fatal(err)
	}
	shell, ok := ev.(cursor.BeforeShellExecution)
	if !ok {
		t.Fatalf("want BeforeShellExecution, got %T", ev)
	}
	if shell.Command != "git push --force" || shell.ConversationID != "c1" {
		t.Fatalf("bad event: %+v", shell)
	}

	out, code, err := cursor.Encode(cursor.EventBeforeShellExecution, cursor.PermissionResultsForTest().Deny("force push blocked"))
	if err != nil {
		t.Fatal(err)
	}
	if code != cursor.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursor.PermissionDenyExit)
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
	ev, err := cursor.DecodeForTest([]byte(cursorStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(cursor.Stop)
	if !ok {
		t.Fatalf("want Stop, got %T", ev)
	}
	if stop.Status != "error" || stop.LoopCount != 1 {
		t.Fatalf("bad stop: %+v", stop)
	}

	out, code, err := cursor.Encode(cursor.EventStop, cursor.StopResultsForTest().FollowUp("retry with fixed creds"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestDecode_AfterFileEdit(t *testing.T) {
	ev, err := cursor.DecodeForTest([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	edit, ok := ev.(cursor.AfterFileEdit)
	if !ok {
		t.Fatalf("want AfterFileEdit, got %T", ev)
	}
	if edit.FilePath != "main.go" || len(edit.Edits) != 1 || edit.Edits[0].OldString != "foo" {
		t.Fatalf("bad edit: %+v", edit)
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	raw := `{"conversation_id":"c1","command":"ls","cwd":"/w"}`
	_, err := cursor.DecodeForTest([]byte(raw))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, cursor.ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "hook_event_name") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := cursor.DecodeForTest([]byte(`{"hook_event_name":"FutureEvent","conversation_id":"c1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":123}`)
	_, err := cursor.DecodeForTest(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, cursor.ErrDecodePayload) {
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
			wantType: cursor.SessionStart{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.SessionStart)
				if e.Model != "gpt" || !e.IsBackgroundAgent {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "sessionEnd",
			raw:      `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`,
			wantType: cursor.SessionEnd{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.SessionEnd)
				if e.Reason != "complete" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeSubmitPrompt",
			raw:      `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`,
			wantType: cursor.BeforeSubmitPrompt{},
			check: func(t *testing.T, ev run.Event) {
				if ev.(cursor.BeforeSubmitPrompt).Prompt != "hello" {
					t.Fatal("bad prompt")
				}
			},
		},
		{
			name:     "preToolUse shell",
			raw:      `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			wantType: cursor.PreToolUse{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.PreToolUse)
				if e.ShellCommand() != "ls" {
					t.Fatalf("ShellCommand=%q", e.ShellCommand())
				}
			},
		},
		{
			name:     "postToolUse",
			raw:      `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`,
			wantType: cursor.PostToolUse{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.PostToolUse)
				if e.ToolOutput != "contents" || e.DurationMillis() != 100 {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "postToolUseFailure",
			raw:      `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`,
			wantType: cursor.PostToolUseFailure{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.PostToolUseFailure)
				if e.ErrorMessage != "timeout" || e.FailureType != "timeout" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeShellExecution",
			raw:      cursorShell,
			wantType: cursor.BeforeShellExecution{},
		},
		{
			name:     "afterShellExecution",
			raw:      `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`,
			wantType: cursor.AfterShellExecution{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.AfterShellExecution)
				if e.Command != "ls" || e.Output != "a\nb" {
					t.Fatalf("event=%+v", e)
				}
			},
		},
		{
			name:     "beforeMCPExecution",
			raw:      `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`,
			wantType: cursor.BeforeMCPExecution{},
			check: func(t *testing.T, ev run.Event) {
				e := ev.(cursor.BeforeMCPExecution)
				if e.URL != "https://mcp.example/mcp" {
					t.Fatalf("URL=%q", e.URL)
				}
			},
		},
		{
			name:     "afterMCPExecution",
			raw:      `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`,
			wantType: cursor.AfterMCPExecution{},
		},
		{
			name:     "beforeReadFile",
			raw:      `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`,
			wantType: cursor.BeforeReadFile{},
		},
		{
			name:     "afterFileEdit",
			raw:      cursorAfterFileEdit,
			wantType: cursor.AfterFileEdit{},
		},
		{
			name:     "subagentStart",
			raw:      `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`,
			wantType: cursor.SubagentStart{},
		},
		{
			name:     "subagentStop",
			raw:      `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`,
			wantType: cursor.SubagentStop{},
		},
		{
			name:     "stop",
			raw:      cursorStop,
			wantType: cursor.Stop{},
		},
		{
			name:     "preCompact",
			raw:      `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`,
			wantType: cursor.PreCompact{},
		},
		{
			name:     "afterAgentResponse",
			raw:      `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`,
			wantType: cursor.AfterAgentResponse{},
		},
		{
			name:     "afterAgentThought",
			raw:      `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`,
			wantType: cursor.AfterAgentThought{},
		},
		{
			name:     "beforeTabFileRead",
			raw:      `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`,
			wantType: cursor.BeforeTabFileRead{},
		},
		{
			name:     "afterTabFileEdit",
			raw:      `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`,
			wantType: cursor.AfterTabFileEdit{},
		},
		{
			name:     "workspaceOpen",
			raw:      `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`,
			wantType: cursor.WorkspaceOpen{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := cursor.DecodeForTest([]byte(tt.raw))
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
	case cursor.SessionStart:
		return "SessionStart"
	case cursor.SessionEnd:
		return "SessionEnd"
	case cursor.BeforeSubmitPrompt:
		return "BeforeSubmitPrompt"
	case cursor.PreToolUse:
		return "PreToolUse"
	case cursor.PostToolUse:
		return "PostToolUse"
	case cursor.PostToolUseFailure:
		return "PostToolUseFailure"
	case cursor.BeforeShellExecution:
		return "BeforeShellExecution"
	case cursor.AfterShellExecution:
		return "AfterShellExecution"
	case cursor.BeforeMCPExecution:
		return "BeforeMCPExecution"
	case cursor.AfterMCPExecution:
		return "AfterMCPExecution"
	case cursor.BeforeReadFile:
		return "BeforeReadFile"
	case cursor.AfterFileEdit:
		return "AfterFileEdit"
	case cursor.SubagentStart:
		return "SubagentStart"
	case cursor.SubagentStop:
		return "SubagentStop"
	case cursor.Stop:
		return "Stop"
	case cursor.PreCompact:
		return "PreCompact"
	case cursor.AfterAgentResponse:
		return "AfterAgentResponse"
	case cursor.AfterAgentThought:
		return "AfterAgentThought"
	case cursor.BeforeTabFileRead:
		return "BeforeTabFileRead"
	case cursor.AfterTabFileEdit:
		return "AfterTabFileEdit"
	case cursor.WorkspaceOpen:
		return "WorkspaceOpen"
	default:
		return "unknown"
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := cursor.Encode(cursor.EventBeforeShellExecution, nil)
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_BeforeSubmitPromptBlock(t *testing.T) {
	out, code, err := cursor.Encode(cursor.EventBeforeSubmitPrompt, cursor.BeforeSubmitPromptResultsForTest().Block("blocked"))
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
	out, code, err := cursor.Encode(cursor.EventSessionStart, cursor.SessionStartResultsForTest().Noop().
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
	out, code, err := cursor.Encode(cursor.EventBeforeTabFileRead, cursor.PermissionResultsForTest().Deny("no tab reads"))
	if err != nil {
		t.Fatal(err)
	}
	if code != cursor.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursor.PermissionDenyExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionUpdatedInputEmptyEventName(t *testing.T) {
	out, code, err := cursor.Encode("", cursor.PermissionResultsForTest().Allow().WithUpdatedInput(map[string]any{"command": "ls"}))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"updated_input"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestMux_Serve_BeforeShellDeny(t *testing.T) {
	run.Reset()
	cursor.OnBeforeShellExecution(func(ctx context.Context, hook run.Hook[cursor.BeforeShellExecution], r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		return r.Deny("blocked"), nil
	})
	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(cursorShell), &stdout, &bytes.Buffer{})
	if code != cursor.PermissionDenyExit {
		t.Fatalf("exit code = %d, want %d", code, cursor.PermissionDenyExit)
	}
	if !strings.Contains(stdout.String(), `"permission":"deny"`) {
		t.Fatalf("bad stdout: %s", stdout.String())
	}
}

func TestToolInput_AsShell(t *testing.T) {
	ev, err := cursor.DecodeForTest([]byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(cursor.PreToolUse)
	input, ok := pre.ToolInput.AsShell()
	if !ok || input.Command != "ls" {
		t.Fatalf("AsShell = %+v, %v", input, ok)
	}
}
