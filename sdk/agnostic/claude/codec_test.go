package claude

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
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

func TestClaudeDecode_InvalidJSONPreservesSentinel(t *testing.T) {
	c := &Codec{}
	_, err := c.Decode([]byte("not json"), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sdkclaude.ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
	if !strings.HasPrefix(err.Error(), "claude: decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestClaudeDecodeEncode_PreToolDeny(t *testing.T) {
	c := &Codec{}
	ev, err := c.Decode([]byte(claudePreToolUse), "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != model.KindPreTool || ev.Session != "abc123" || ev.Cwd != "/home/user/proj" {
		t.Fatalf("bad event: %+v", ev)
	}
	if ev.Tool == nil || ev.Tool.Name != model.ToolBash || ev.Tool.Native != "Bash" || ev.Tool.Shell != "rm -rf /tmp/build" {
		t.Fatalf("bad tool: %+v", ev.Tool)
	}
	if !bytes.Equal(ev.Raw, []byte(claudePreToolUse)) {
		t.Fatal("Raw not preserved")
	}

	out, code, err := c.Encode(ev, model.Deny("destructive command"))
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
	c := &Codec{}
	ev, err := c.Decode([]byte(claudePreToolUse), "")
	if err != nil {
		t.Fatal(err)
	}

	out, code, err := c.Encode(ev, model.Result{UpdatedInput: map[string]any{"command": "ls -la"}})
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
	c := &Codec{}
	stopEv := &model.Event{Agent: model.Claude, Kind: model.KindStop, Name: "Stop"}
	out, code, err := c.Encode(stopEv, model.Result{FollowUp: "run the tests"})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad stop output: %s", out)
	}
}

func TestClaudeEncode_SubagentStopFollowUp(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Claude, Kind: model.KindSubagentStop, Name: "SubagentStop"}
	out, code, err := c.Encode(ev, model.Result{FollowUp: "finish the review"})
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "finish the review") {
		t.Fatalf("bad subagent stop output: %s", out)
	}
}

func TestClaudeEncode_SessionStartContext(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Claude, Kind: model.KindSessionStart, Name: "SessionStart"}
	out, code, err := c.Encode(ev, model.Context("boot notes"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"additionalContext":"boot notes"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeEncode_ZeroResult(t *testing.T) {
	c := &Codec{}
	ev, _ := c.Decode([]byte(claudePreToolUse), "")
	out, code, err := c.Encode(ev, model.Result{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero result should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestClaudeDecode_UnknownEvent(t *testing.T) {
	raw := []byte(`{"session_id":"s1","hook_event_name":"Setup","cwd":"/w"}`)
	c := &Codec{}
	ev, err := c.Decode(raw, "")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Kind != model.KindOther || ev.Name != "Setup" {
		t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
	}
	if !bytes.Equal(ev.Raw, raw) {
		t.Fatal("Raw not preserved")
	}
}

func TestClaudeDecode_Matrix(t *testing.T) {
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
	c := &Codec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := c.Decode([]byte(tt.raw), "")
			if err != nil {
				t.Fatal(err)
			}
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

func TestClaudeEncode_NilEvent(t *testing.T) {
	c := &Codec{}
	_, _, err := c.Encode(nil, model.Deny("nope"))
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
	c := &Codec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.raw)
			ev, err := c.Decode(raw, "")
			if err != nil {
				t.Fatal(err)
			}
			if ev.Kind != model.KindOther || ev.Name != tt.name {
				t.Fatalf("model.Kind=%v Name=%q", ev.Kind, ev.Name)
			}
			if !bytes.Equal(ev.Raw, raw) {
				t.Fatal("Raw not preserved")
			}
		})
	}
}

func TestClaudeDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"session_id":"s","hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"ls"}}`)
	c := &Codec{}
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

func TestClaudeEncode_PostToolUpdatedOutput(t *testing.T) {
	c := &Codec{}
	ev := &model.Event{Agent: model.Claude, Kind: model.KindPostTool, Name: "PostToolUse"}
	text := "rewritten"
	out, code, err := c.Encode(ev, model.Result{UpdatedOutput: &text})
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"updatedToolOutput":"rewritten"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestClaudeDecode_InvalidJSON(t *testing.T) {
	c := &Codec{}
	_, err := c.Decode([]byte("not json"), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "claude:") {
		t.Fatalf("error = %v", err)
	}
}
