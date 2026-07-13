package agenthooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cursorhook"
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

func TestCursorDecodeEncode_ShellGuard(t *testing.T) {
	c := &CursorCodec{}
	ev, err := c.Decode([]byte(cursorShell), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != KindPreTool || ev.Name != "beforeShellExecution" || ev.Session != "c1" {
		t.Fatalf("bad event: %+v", ev)
	}
	if ev.Tool == nil || ev.Tool.Name != ToolBash || ev.Tool.Shell != "git push --force" {
		t.Fatalf("bad tool: %+v", ev.Tool)
	}
	if !bytes.Equal(ev.Raw, []byte(cursorShell)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := c.Encode(ev, Deny("force push blocked"))
	if err != nil {
		t.Fatal(err)
	}
	if code != CursorWarnExit {
		t.Fatalf("exit code = %d, want %d", code, CursorWarnExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "deny" || got["agent_message"] != "force push blocked" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCursorDecodeEncode_StopFollowUp(t *testing.T) {
	c := &CursorCodec{}
	ev, err := c.Decode([]byte(cursorStop), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != KindStop || ev.Turn == nil || ev.Turn.Status != "error" || ev.Turn.LoopCount != 1 {
		t.Fatalf("bad stop: %+v turn=%+v", ev, ev.Turn)
	}

	out, code, err := c.Encode(ev, Result{FollowUp: "retry with fixed creds"})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"retry with fixed creds"`) {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestCursorDecode_AfterFileEdit(t *testing.T) {
	c := &CursorCodec{}
	ev, err := c.Decode([]byte(cursorAfterFileEdit), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != KindPostTool || ev.Name != "afterFileEdit" {
		t.Fatalf("Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if ev.Tool == nil || ev.Tool.Name != ToolEdit || ev.Tool.Native != "afterFileEdit" {
		t.Fatalf("Tool=%+v", ev.Tool)
	}
	if !bytes.Equal(ev.Raw, []byte(cursorAfterFileEdit)) {
		t.Fatal("Raw not preserved")
	}

	var input struct {
		FilePath string `json:"file_path"`
		Edits    []struct {
			OldString string `json:"old_string"`
			NewString string `json:"new_string"`
		} `json:"edits"`
	}
	if err := json.Unmarshal(ev.Tool.Input, &input); err != nil {
		t.Fatal(err)
	}
	if input.FilePath != "main.go" || len(input.Edits) != 1 || input.Edits[0].OldString != "foo" {
		t.Fatalf("bad input: %+v", input)
	}
	if ev.Result == nil || len(ev.Result.Raw) == 0 {
		t.Fatalf("Result.Raw missing: %+v", ev.Result)
	}
}

func TestCursorDecode_Matrix(t *testing.T) {
	c := &CursorCodec{}
	tests := []struct {
		name      string
		raw       string
		eventHint string
		kind      Kind
		check     func(t *testing.T, ev *Event)
	}{
		{
			name: "sessionStart",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`,
			kind: KindSessionStart,
			check: func(t *testing.T, ev *Event) {
				if ev.Life == nil || ev.Life.Model != "gpt" || !ev.Life.Background {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "sessionEnd",
			raw:  `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`,
			kind: KindSessionEnd,
			check: func(t *testing.T, ev *Event) {
				if ev.Life == nil || ev.Life.Reason != "complete" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "beforeSubmitPrompt",
			raw:  `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`,
			kind: KindUserPrompt,
			check: func(t *testing.T, ev *Event) {
				if ev.Prompt != "hello" {
					t.Fatalf("Prompt=%q", ev.Prompt)
				}
			},
		},
		{
			name: "preToolUse shell",
			raw:  `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			kind: KindPreTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Name != ToolBash || ev.Tool.Shell != "ls" {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "postToolUse",
			raw:  `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`,
			kind: KindPostTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Result == nil || ev.Result.Text != "contents" || ev.Result.DurationMs != 100 {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "postToolUseFailure",
			raw:  `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`,
			kind: KindPostToolFailure,
			check: func(t *testing.T, ev *Event) {
				if ev.Result == nil || ev.Result.Error != "timeout" || ev.Result.FailureType != "timeout" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "beforeShellExecution",
			raw:  cursorShell,
			kind: KindPreTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Name != "beforeShellExecution" || ev.Tool.Shell != "git push --force" {
					t.Fatalf("event=%+v tool=%+v", ev, ev.Tool)
				}
			},
		},
		{
			name: "afterShellExecution",
			raw:  `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`,
			kind: KindPostTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Name != ToolBash || ev.Result.Text != "a\nb" {
					t.Fatalf("tool=%+v result=%+v", ev.Tool, ev.Result)
				}
			},
		},
		{
			name: "beforeMCPExecution",
			raw:  `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}"}`,
			kind: KindPreTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || !ev.Tool.MCP || ev.Tool.Name != "MCP:browser_navigate" {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "afterMCPExecution",
			raw:  `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`,
			kind: KindPostTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || !ev.Tool.MCP || ev.Result.Text != "{}" {
					t.Fatalf("tool=%+v result=%+v", ev.Tool, ev.Result)
				}
			},
		},
		{
			name: "beforeReadFile",
			raw:  `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`,
			kind: KindPreTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Name != ToolRead {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
				var input struct {
					FilePath string `json:"file_path"`
					Content  string `json:"content"`
				}
				if err := json.Unmarshal(ev.Tool.Input, &input); err != nil {
					t.Fatal(err)
				}
				if input.FilePath != "a.go" || input.Content != "package main" {
					t.Fatalf("input=%+v", input)
				}
			},
		},
		{
			name: "afterFileEdit",
			raw:  cursorAfterFileEdit,
			kind: KindPostTool,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool == nil || ev.Tool.Name != ToolEdit {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "subagentStart",
			raw:  `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`,
			kind: KindSubagentStart,
			check: func(t *testing.T, ev *Event) {
				if ev.Subagent == nil || ev.Subagent.ID != "sa1" || ev.Subagent.Type != "explore" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "subagentStop",
			raw:  `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`,
			kind: KindSubagentStop,
			check: func(t *testing.T, ev *Event) {
				if ev.Subagent == nil || ev.Subagent.LoopCount != 2 {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "stop",
			raw:  cursorStop,
			kind: KindStop,
			check: func(t *testing.T, ev *Event) {
				if ev.Turn == nil || ev.Turn.LoopCount != 1 {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "preCompact",
			raw:  `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`,
			kind: KindPreCompact,
			check: func(t *testing.T, ev *Event) {
				if ev.Compact == nil || ev.Compact.Trigger != "auto" {
					t.Fatalf("Compact=%+v", ev.Compact)
				}
			},
		},
		{
			name: "afterAgentResponse",
			raw:  `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`,
			kind: KindOther,
			check: func(t *testing.T, ev *Event) {
				if ev.Note == nil || ev.Note.Message != "done" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "afterAgentThought",
			raw:  `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`,
			kind: KindOther,
			check: func(t *testing.T, ev *Event) {
				if ev.Note == nil || ev.Note.Type != "thought" || ev.Note.Message != "thinking" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "beforeTabFileRead",
			raw:  `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`,
			kind: KindOther,
			check: func(t *testing.T, ev *Event) {
				if ev.Tool != nil {
					t.Fatalf("tab event should not fold Tool, got %+v", ev.Tool)
				}
			},
		},
		{
			name: "afterTabFileEdit",
			raw:  `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`,
			kind: KindOther,
		},
		{
			name: "workspaceOpen",
			raw:  `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`,
			kind: KindOther,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := c.Decode([]byte(tt.raw), tt.eventHint)
			if err != nil {
				t.Fatal(err)
			}
			if ev.Kind != tt.kind {
				t.Fatalf("Kind=%v want %v", ev.Kind, tt.kind)
			}
			if ev.Name == "" {
				t.Fatal("Name empty")
			}
			if len(ev.Raw) == 0 {
				t.Fatal("Raw empty")
			}
			if tt.check != nil {
				tt.check(t, ev)
			}
		})
	}
}

func TestCursorDecode_RequiresEventHint(t *testing.T) {
	raw := `{"conversation_id":"c1","command":"ls","cwd":"/w"}`
	c := &CursorCodec{}
	_, err := c.Decode([]byte(raw), "")
	if err == nil {
		t.Fatal("expected error without eventHint")
	}
	if !errors.Is(err, cursorhook.ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "eventHint") {
		t.Fatalf("error = %v", err)
	}
}

func TestCursorEncode_ZeroResult(t *testing.T) {
	c := &CursorCodec{}
	ev := &Event{Agent: Cursor, Kind: KindPreTool, Name: "beforeShellExecution"}
	out, code, err := c.Encode(ev, Result{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestCursorEncode_BeforeSubmitPromptBlock(t *testing.T) {
	c := &CursorCodec{}
	ev := &Event{Agent: Cursor, Kind: KindUserPrompt, Name: "beforeSubmitPrompt"}
	out, code, err := c.Encode(ev, Result{BlockPrompt: true, UserMessage: "blocked"})
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

func TestCursorEncode_SessionStartEnv(t *testing.T) {
	c := &CursorCodec{}
	ev := &Event{Agent: Cursor, Kind: KindSessionStart, Name: "sessionStart"}
	out, code, err := c.Encode(ev, Result{Env: map[string]string{"K": "V"}, Context: "ctx"})
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

func TestCursorEncode_TabFileReadDeny(t *testing.T) {
	c := &CursorCodec{}
	ev := &Event{Agent: Cursor, Kind: KindOther, Name: "beforeTabFileRead"}
	out, code, err := c.Encode(ev, Deny("no tab reads"))
	if err != nil {
		t.Fatal(err)
	}
	if code != CursorWarnExit {
		t.Fatalf("exit code = %d, want %d", code, CursorWarnExit)
	}
	if !strings.Contains(string(out), `"permission":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestCursorEncode_NilEvent(t *testing.T) {
	c := &CursorCodec{}
	_, _, err := c.Encode(nil, Deny("nope"))
	if err == nil {
		t.Fatal("expected error for nil event")
	}
}
