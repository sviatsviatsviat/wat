package cursor

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
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

func mapRaw(t *testing.T, raw []byte, eventHint string) *model.Event {
	t.Helper()
	var opts []sdkcursor.Option
	if eventHint != "" {
		opts = append(opts, sdkcursor.WithEvent(eventHint))
	}
	native, err := sdkcursor.Decode(raw, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return MapEvent(native, raw)
}

func TestCursorMapEvent_AfterFileEdit(t *testing.T) {
	ev := mapRaw(t, []byte(cursorAfterFileEdit), "")
	if ev.Kind != model.KindPostTool || ev.Name != "afterFileEdit" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if ev.Tool == nil || ev.Tool.Name != model.ToolEdit || ev.Tool.Native != "afterFileEdit" {
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
	if err := json.Unmarshal(ev.Tool.Input.Raw(), &input); err != nil {
		t.Fatal(err)
	}
	if input.FilePath != "main.go" || len(input.Edits) != 1 || input.Edits[0].OldString != "foo" {
		t.Fatalf("bad input: %+v", input)
	}
	if ev.Result == nil || len(ev.Result.Raw) == 0 {
		t.Fatalf("Result.Raw missing: %+v", ev.Result)
	}
}

func TestCursorMapEvent_Matrix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		eventHint string
		kind      model.Kind
		check     func(t *testing.T, ev *model.Event)
	}{
		{
			name: "sessionStart",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`,
			kind: model.KindSessionStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Model != "gpt" || !ev.Life.Background {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "sessionEnd",
			raw:  `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`,
			kind: model.KindSessionEnd,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Reason != "complete" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "beforeSubmitPrompt",
			raw:  `{"hook_event_name":"beforeSubmitPrompt","conversation_id":"c1","prompt":"hello"}`,
			kind: model.KindUserPrompt,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Prompt != "hello" {
					t.Fatalf("Prompt=%q", ev.Prompt)
				}
			},
		},
		{
			name: "preToolUse shell",
			raw:  `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			kind: model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Shell != "ls" {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "postToolUse",
			raw:  `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`,
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Result == nil || ev.Result.Text != "contents" || ev.Result.DurationMs != 100 {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "postToolUseFailure",
			raw:  `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`,
			kind: model.KindPostToolFailure,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Result == nil || ev.Result.Error != "timeout" || ev.Result.FailureType != "timeout" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "beforeShellExecution",
			raw:  cursorShell,
			kind: model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Name != "beforeShellExecution" || ev.Tool.Shell != "git push --force" {
					t.Fatalf("event=%+v tool=%+v", ev, ev.Tool)
				}
			},
		},
		{
			name: "afterShellExecution",
			raw:  `{"hook_event_name":"afterShellExecution","conversation_id":"c1","command":"ls","output":"a\nb","duration":10}`,
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Result.Text != "a\nb" {
					t.Fatalf("tool=%+v result=%+v", ev.Tool, ev.Result)
				}
			},
		},
		{
			name: "beforeMCPExecution",
			raw:  `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}"}`,
			kind: model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || !ev.Tool.MCP || ev.Tool.Name != "MCP:browser_navigate" {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "afterMCPExecution",
			raw:  `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration_ms":5}`,
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || !ev.Tool.MCP || ev.Result == nil || string(ev.Result.Raw) != "{}" {
					t.Fatalf("tool=%+v result=%+v", ev.Tool, ev.Result)
				}
			},
		},
		{
			name: "beforeReadFile",
			raw:  `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`,
			kind: model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolRead {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
				var input struct {
					FilePath string `json:"file_path"`
					Content  string `json:"content"`
				}
				if err := json.Unmarshal(ev.Tool.Input.Raw(), &input); err != nil {
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
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolEdit {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "subagentStart",
			raw:  `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`,
			kind: model.KindSubagentStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.ID != "sa1" || ev.Subagent.Type != "explore" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "subagentStop",
			raw:  `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`,
			kind: model.KindSubagentStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.LoopCount != 2 {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "stop",
			raw:  cursorStop,
			kind: model.KindStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Turn == nil || ev.Turn.LoopCount != 1 {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "preCompact",
			raw:  `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"auto"}`,
			kind: model.KindPreCompact,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Compact == nil || ev.Compact.Trigger != "auto" {
					t.Fatalf("Compact=%+v", ev.Compact)
				}
			},
		},
		{
			name: "afterAgentResponse",
			raw:  `{"hook_event_name":"afterAgentResponse","conversation_id":"c1","text":"done"}`,
			kind: model.KindOther,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Note == nil || ev.Note.Message != "done" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "afterAgentThought",
			raw:  `{"hook_event_name":"afterAgentThought","conversation_id":"c1","text":"thinking"}`,
			kind: model.KindOther,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Note == nil || ev.Note.Type != "thought" || ev.Note.Message != "thinking" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "beforeTabFileRead",
			raw:  `{"hook_event_name":"beforeTabFileRead","conversation_id":"c1","file_path":"x.go","content":"x"}`,
			kind: model.KindOther,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool != nil {
					t.Fatalf("tab event should not fold Tool, got %+v", ev.Tool)
				}
			},
		},
		{
			name: "afterTabFileEdit",
			raw:  `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`,
			kind: model.KindOther,
		},
		{
			name: "workspaceOpen",
			raw:  `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`,
			kind: model.KindOther,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := mapRaw(t, []byte(tt.raw), tt.eventHint)
			if ev.Kind != tt.kind {
				t.Fatalf("model.Kind=%v want %v", ev.Kind, tt.kind)
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
