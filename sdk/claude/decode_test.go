package claude_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
)

func TestDecode_PreToolUse(t *testing.T) {
	ev, err := claude.Decode([]byte(preToolUsePayload))
	if err != nil {
		t.Fatal(err)
	}
	pre, ok := ev.(claude.PreToolUse)
	if !ok || pre.ToolName != "Bash" || pre.SessionID != "abc123" {
		t.Fatalf("bad event: %+v", ev)
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
			ev, err := claude.Decode([]byte(tt.raw))
			if err != nil {
				t.Fatal(err)
			}
			if ev.EventName() != tt.name {
				t.Fatalf("EventName() = %q, want %q", ev.EventName(), tt.name)
			}
		})
	}
}

func TestDecode_UnknownEventPanics(t *testing.T) {
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("expected panic")
		}
		msg, ok := got.(string)
		if !ok || !strings.Contains(msg, "unknown hook event") {
			t.Fatalf("recover = %#v", got)
		}
	}()
	_, _ = claude.Decode([]byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`))
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := claude.Decode([]byte(`{"session_id":"s1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, claude.ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	t.Run("envelope", func(t *testing.T) {
		_, err := claude.Decode([]byte("not json"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claude.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("typed event", func(t *testing.T) {
		raw := []byte(`{"session_id":"s1","hook_event_name":"PreToolUse","cwd":"/w","tool_name":123}`)
		_, err := claude.Decode(raw)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, claude.ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})
}
