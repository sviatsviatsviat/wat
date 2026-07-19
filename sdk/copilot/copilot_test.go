package copilot

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

const copilotPreToolUse = `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "bash",
  "tool_input": {"command": "rm -rf /"}
}`
const copilotVSCodeStop = `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "transcript_path": "/tmp/t",
  "stop_reason": "end_turn"
}`

func TestDecodeEncode_PreToolDeny(t *testing.T) {
	ev, err := codec.Decode([]byte(copilotPreToolUse))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(PreToolUse)
	if !ok || pre.SessionID != "s1" || pre.Cwd != "/w" {
		t.Fatalf("bad event: %+v", ev)
	}
	if pre.NativeToolName() != "bash" || pre.ShellCommand() != "rm -rf /" {
		t.Fatalf("bad tool: name=%q shell=%q", pre.NativeToolName(), pre.ShellCommand())
	}

	out, code, err := encode(EventPreToolUse, preToolResults{}.Deny("destructive command"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Decision string `json:"permission_decision"`
		Reason   string `json:"permission_decision_reason"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Decision != "deny" || got.Reason != "destructive command" {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := codec.Decode([]byte(`{"session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := codec.Decode([]byte(`{"hook_event_name":"FutureEvent","session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s1","cwd":"/w","tool_name":123}`)
	_, err := codec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestDecode_VSCodeStop(t *testing.T) {
	ev, err := codec.Decode([]byte(copilotVSCodeStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(AgentStop)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if stop.EventName() != EventAgentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if stop.HookEventName != "Stop" {
		t.Fatalf("HookEventName=%q", stop.HookEventName)
	}
	if stop.Reason() != "end_turn" {
		t.Fatalf("Reason=%q", stop.Reason())
	}
	if stop.TranscriptPath != "/tmp/t" {
		t.Fatalf("TranscriptPath=%q", stop.TranscriptPath)
	}
}

func TestDecode_VSCodePreToolBash(t *testing.T) {
	raw := `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "Bash",
  "tool_input": {"command": "ls -la"}
}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(PreToolUse)
	if !ok || pre.NativeToolName() != "Bash" || pre.ShellCommand() != "ls -la" {
		t.Fatalf("PreToolUse=%+v", ev)
	}
}

func TestDecode_Notification(t *testing.T) {
	raw := `{
  "hook_event_name": "Notification",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "message": "shell done",
  "title": "Shell completed",
  "notification_type": "shell_completed"
}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	note, ok := ev.(Notification)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if note.EventName() != EventNotification {
		t.Fatalf("EventName=%q", note.EventName())
	}
	if note.NotificationType != "shell_completed" || note.Message != "shell done" {
		t.Fatalf("Notification=%+v", note)
	}
}

func TestDecode_Matrix(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		eventName string
		check     func(t *testing.T, ev run.Event)
	}{
		{
			name:      "sessionStart",
			raw:       `{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`,
			eventName: EventSessionStart,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(SessionStart)
				if !ok || e.Source != "new" || e.InitialPrompt() != "go" {
					t.Fatalf("SessionStart=%+v", ev)
				}
			},
		},
		{
			name:      "sessionEnd",
			raw:       `{"hook_event_name":"SessionEnd","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","reason":"complete"}`,
			eventName: EventSessionEnd,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(SessionEnd)
				if !ok || e.Reason != "complete" {
					t.Fatalf("SessionEnd=%+v", ev)
				}
			},
		},
		{
			name:      "userPromptSubmitted",
			raw:       `{"hook_event_name":"UserPromptSubmit","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","prompt":"hello"}`,
			eventName: EventUserPromptSubmitted,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(UserPromptSubmitted)
				if !ok || e.Prompt != "hello" {
					t.Fatalf("UserPromptSubmitted=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUse",
			raw:       `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"view","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"contents"}}`,
			eventName: EventPostToolUse,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(PostToolUse)
				if !ok || e.ResultText() != "contents" {
					t.Fatalf("PostToolUse=%+v", ev)
				}
			},
		},
		{
			name:      "postToolUseFailure",
			raw:       `{"hook_event_name":"PostToolUseFailure","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{},"error":"timeout"}`,
			eventName: EventPostToolUseFailure,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(PostToolUseFailure)
				if !ok || e.ErrorMessage() != "timeout" {
					t.Fatalf("PostToolUseFailure=%+v", ev)
				}
			},
		},
		{
			name:      "permissionRequest",
			raw:       `{"hook_event_name":"PermissionRequest","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"create","tool_input":{"path":"a.txt"}}`,
			eventName: EventPermissionRequest,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(PermissionRequest)
				if !ok || e.NativeToolName() != "create" {
					t.Fatalf("PermissionRequest=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStart",
			raw:       `{"hook_event_name":"SubagentStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"explore","agent_display_name":"Explore","agent_description":"search codebase"}`,
			eventName: EventSubagentStart,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(SubagentStart)
				if !ok || e.Name() != "explore" || e.DisplayName() != "Explore" {
					t.Fatalf("SubagentStart=%+v", ev)
				}
			},
		},
		{
			name:      "subagentStop explicit",
			raw:       `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","agent_display_name":"Task","stop_reason":"end_turn"}`,
			eventName: EventSubagentStop,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(SubagentStop)
				if !ok || e.Name() != "task" || e.Reason() != "end_turn" {
					t.Fatalf("SubagentStop=%+v", ev)
				}
			},
		},
		{
			name:      "agentStop",
			raw:       `{"hook_event_name":"Stop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","stop_reason":"end_turn"}`,
			eventName: EventAgentStop,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(AgentStop)
				if !ok || e.Reason() != "end_turn" || e.IsSubagent() {
					t.Fatalf("AgentStop=%+v", ev)
				}
			},
		},
		{
			name:      "agentStop with agent scope",
			raw:       `{"hook_event_name":"Stop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`,
			eventName: EventAgentStop,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(AgentStop)
				if !ok || !e.IsSubagent() || e.Name() != "task" {
					t.Fatalf("AgentStop=%+v", ev)
				}
			},
		},
		{
			name:      "preCompact",
			raw:       `{"hook_event_name":"PreCompact","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","trigger":"auto","custom_instructions":"keep"}`,
			eventName: EventPreCompact,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(PreCompact)
				if !ok || e.Instructions() != "keep" {
					t.Fatalf("PreCompact=%+v", ev)
				}
			},
		},
		{
			name:      "errorOccurred",
			raw:       `{"hook_event_name":"ErrorOccurred","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","error":{"message":"slow down","name":"RateLimit"},"error_context":"model_call","recoverable":true}`,
			eventName: EventErrorOccurred,
			check: func(t *testing.T, ev run.Event) {
				e, ok := ev.(ErrorOccurred)
				if !ok {
					t.Fatalf("ErrorOccurred=%+v", ev)
				}
				detail, ok := e.Detail()
				if !ok || detail.Name != "RateLimit" || detail.Message != "slow down" {
					t.Fatalf("Detail=%+v ok=%v", detail, ok)
				}
				if e.Recoverable == nil || !*e.Recoverable {
					t.Fatalf("Recoverable=%v", e.Recoverable)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := codec.Decode([]byte(tt.raw))
			if err != nil {
				t.Fatal(err)
			}
			if ev.EventName() != tt.eventName {
				t.Fatalf("EventName=%q, want %q", ev.EventName(), tt.eventName)
			}
			tt.check(t, ev)
		})
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	_, err := codec.Decode([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "decode payload") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_ToolInputNotAliased(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{"command":"ls"}}`)
	ev, err := codec.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	got := pre.ToolInput.Raw()
	got[0] = 'X'
	if bytes.Equal(pre.ToolInput.Raw(), got) {
		t.Fatal("ToolInput.Raw() did not return a defensive copy")
	}
}

func TestEncode_PreToolAllowModifiedArgs(t *testing.T) {
	out, code, err := encode(EventPreToolUse, preToolResults{}.Allow().WithModifiedArgs(map[string]any{"command": "echo safe"}))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"permission_decision":"allow"`) || !strings.Contains(string(out), `"modified_args"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolUpdatedOutput(t *testing.T) {
	out, code, err := encode(EventPostToolUse, postToolResults{}.Context("extra guidance").WithModifiedResult("rewritten"))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	s := string(out)
	if !strings.Contains(s, `"text_result_for_llm":"rewritten"`) || !strings.Contains(s, `"additional_context":"extra guidance"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopFollowUp(t *testing.T) {
	out, code, err := encode(EventAgentStop, stopResults{}.FollowUp("run the tests"))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestDenyInterrupt(t *testing.T) {
	out, code, err := encode(EventPermissionRequest, permissionRequestResults{}.Deny("blocked").WithInterrupt(true))
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) || !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PermissionRequestAsk(t *testing.T) {
	out, code, err := encode(EventPermissionRequest, permissionRequestResults{}.Ask("needs user confirmation"))
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"behavior":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_PostToolFailureContext(t *testing.T) {
	out, code, err := encode(EventPostToolUseFailure, postToolFailureResults{}.Context("retry with smaller input"))
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := encode(EventSessionStart, sessionStartResults{}.Context("project uses go test ./..."))
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additional_context"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := encode(EventPreToolUse, nil)
	if err != nil || code != 0 || out != nil {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, _, err := encode(EventPostToolUse, preToolResults{}.Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestErrorOccurred_DetailNull(t *testing.T) {
	e := ErrorOccurred{Error: json.RawMessage("null")}
	if _, ok := e.Detail(); ok {
		t.Fatal("JSON null error payload should be absent")
	}
}

func TestToolInput_AsBash(t *testing.T) {
	ev, err := codec.Decode([]byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"Bash","tool_input":{"command":"ls -la"}}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(PreToolUse)
	input, ok := pre.Input().AsBash()
	if !ok || input.Command != "ls -la" {
		t.Fatalf("AsBash = %+v, %v", input, ok)
	}
}

func TestDecode_VSCodeStopWithAgentScope(t *testing.T) {
	raw := `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "agent_name": "task",
  "stop_reason": "end_turn"
}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(AgentStop)
	if !ok {
		t.Fatalf("got %T, want AgentStop", ev)
	}
	if stop.EventName() != EventAgentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if !stop.IsSubagent() || stop.Name() != "task" {
		t.Fatalf("IsSubagent/Name = %v %q", stop.IsSubagent(), stop.Name())
	}
	if stop.HookEventName != "Stop" {
		t.Fatalf("HookEventName=%q", stop.HookEventName)
	}
}

func TestDecode_PostToolUseResultRaw(t *testing.T) {
	raw := `{"hook_event_name":"PostToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"view","tool_input":{},"tool_result":{"result_type":"success","text_result_for_llm":"contents"}}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	post := ev.(PostToolUse)
	got := string(post.ResultRaw())
	if !strings.Contains(got, "text_result_for_llm") || !strings.Contains(got, "contents") {
		t.Fatalf("ResultRaw=%s", got)
	}
}
