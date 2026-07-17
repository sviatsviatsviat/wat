package claude

import (
	"bytes"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
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
// Preserved only via Raw (model.KindOther):
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

func mapRaw(t *testing.T, raw []byte) *model.Event {
	t.Helper()
	native, err := sdkclaude.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	return MapEvent(native, raw)
}

func TestClaudeMapEvent_UnknownEvent(t *testing.T) {
	raw := []byte(`{"session_id":"s1","hook_event_name":"Setup","cwd":"/w"}`)
	ev := mapRaw(t, raw)
	if ev.Kind != model.KindOther || ev.Name != "Setup" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if !bytes.Equal(ev.Raw, raw) {
		t.Fatal("Raw not preserved")
	}
}

func TestClaudeMapEvent_Matrix(t *testing.T) {
	tests := []struct {
		name  string
		raw   string
		kind  model.Kind
		check func(t *testing.T, ev *model.Event)
	}{
		{
			name: "SessionStart",
			raw:  `{"session_id":"s","hook_event_name":"SessionStart","source":"startup","model":"claude-3"}`,
			kind: model.KindSessionStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Source != "startup" || ev.Life.Model != "claude-3" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "SessionEnd",
			raw:  `{"session_id":"s","hook_event_name":"SessionEnd","reason":"clear"}`,
			kind: model.KindSessionEnd,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Life == nil || ev.Life.Reason != "clear" {
					t.Fatalf("Life=%+v", ev.Life)
				}
			},
		},
		{
			name: "UserPromptSubmit",
			raw:  `{"session_id":"s","hook_event_name":"UserPromptSubmit","prompt":"hello"}`,
			kind: model.KindUserPrompt,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Prompt != "hello" {
					t.Fatalf("Prompt=%q", ev.Prompt)
				}
			},
		},
		{
			name: "PreToolUse",
			raw:  `{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"},"tool_use_id":"t1"}`,
			kind: model.KindPreTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Shell != "ls" || ev.Tool.Name != model.ToolBash {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "PostToolUse",
			raw:  `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`,
			kind: model.KindPostTool,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Result == nil || ev.Result.Text != "file contents" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "PostToolUseFailure",
			raw:  `{"session_id":"s","hook_event_name":"PostToolUseFailure","tool_name":"Bash","error":"timeout"}`,
			kind: model.KindPostToolFailure,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Result == nil || ev.Result.Error != "timeout" {
					t.Fatalf("Result=%+v", ev.Result)
				}
			},
		},
		{
			name: "PermissionRequest",
			raw:  `{"session_id":"s","hook_event_name":"PermissionRequest","tool_name":"Write","tool_use_id":"t2"}`,
			kind: model.KindPermissionRequest,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Tool == nil || ev.Tool.Name != model.ToolWrite {
					t.Fatalf("Tool=%+v", ev.Tool)
				}
			},
		},
		{
			name: "SubagentStart",
			raw:  `{"session_id":"s","hook_event_name":"SubagentStart","agent_id":"a1","agent_type":"reviewer"}`,
			kind: model.KindSubagentStart,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Subagent == nil || ev.Subagent.ID != "a1" || ev.Subagent.Type != "reviewer" {
					t.Fatalf("Subagent=%+v", ev.Subagent)
				}
			},
		},
		{
			name: "SubagentStop",
			raw:  `{"session_id":"s","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"worker","stop_hook_active":true,"last_assistant_message":"done"}`,
			kind: model.KindSubagentStop,
			check: func(t *testing.T, ev *model.Event) {
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
			kind: model.KindStop,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Turn == nil || ev.Turn.LastAssistantMessage != "bye" {
					t.Fatalf("Turn=%+v", ev.Turn)
				}
			},
		},
		{
			name: "PreCompact",
			raw:  `{"session_id":"s","hook_event_name":"PreCompact","trigger":"manual","custom_instructions":"keep tests"}`,
			kind: model.KindPreCompact,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Compact == nil || ev.Compact.Trigger != "manual" || ev.Compact.CustomInstructions != "keep tests" {
					t.Fatalf("Compact=%+v", ev.Compact)
				}
			},
		},
		{
			name: "Notification",
			raw:  `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`,
			kind: model.KindNotification,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Note == nil || ev.Note.Type != "idle_prompt" || ev.Note.Message != "waiting" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
		{
			name: "StopFailure",
			raw:  `{"session_id":"s","hook_event_name":"StopFailure","error_type":"rate_limit","message":"slow down"}`,
			kind: model.KindAgentError,
			check: func(t *testing.T, ev *model.Event) {
				if ev.Note == nil || ev.Note.Type != "rate_limit" || ev.Note.Message != "slow down" {
					t.Fatalf("Note=%+v", ev.Note)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := mapRaw(t, []byte(tt.raw))
			if ev.Kind != tt.kind {
				t.Fatalf("model.Kind=%v, want %v", ev.Kind, tt.kind)
			}
			if ev.Name != tt.name {
				t.Fatalf("Name=%q", ev.Name)
			}
			tt.check(t, ev)
		})
	}
}

func TestClaudeMapEvent_LongTailKindOther(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.raw)
			ev := mapRaw(t, raw)
			if ev.Kind != model.KindOther || ev.Name != tt.name {
				t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
			}
			if !bytes.Equal(ev.Raw, raw) {
				t.Fatal("Raw not preserved")
			}
		})
	}
}

func TestClaudeMapEvent_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"}}`)
	ev := mapRaw(t, raw)
	rawCopy := bytes.Clone(ev.Raw)
	got := ev.Tool.Input.Raw()
	got[0] = 'X'
	if !bytes.Equal(ev.Raw, rawCopy) {
		t.Fatal("mutating Tool.Input.Raw() copy affected Event.Raw")
	}
	if bytes.Equal(ev.Tool.Input.Raw(), got) {
		t.Fatal("Tool.Input.Raw() did not return a defensive copy")
	}
}
