package claude

import (
	"errors"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte(`{"session_id":"s1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	t.Run("envelope", func(t *testing.T) {
		_, err := runtime.Codec.Decode([]byte("not json"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
	})

	t.Run("typed event", func(t *testing.T) {
		raw := []byte(`{"session_id":"s1","hook_event_name":"PreToolUse","cwd":"/w","tool_name":123}`)
		_, err := runtime.Codec.Decode(raw)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
	})
}

func TestEventNames(t *testing.T) {
	cases := []hookkit.Event{
		SessionStart{},
		Setup{},
		SessionEnd{},
		UserPromptSubmit{},
		UserPromptExpansion{},
		PreToolUse{},
		PostToolUse{},
		PostToolUseFailure{},
		PostToolBatch{},
		PermissionRequest{},
		PermissionDenied{},
		SubagentStart{},
		SubagentStop{},
		TaskCreated{},
		TaskCompleted{},
		Stop{},
		StopFailure{},
		TeammateIdle{},
		Notification{},
		MessageDisplay{},
		InstructionsLoaded{},
		ConfigChange{},
		CwdChanged{},
		FileChanged{},
		WorktreeCreate{},
		WorktreeRemove{},
		PreCompact{},
		PostCompact{},
		Elicitation{},
		ElicitationResult{},
	}
	for _, ev := range cases {
		t.Run(ev.EventName(), func(t *testing.T) {
			if ev.EventName() == "" {
				t.Fatal("EventName() empty")
			}
		})
	}
}

func TestEventNameFromRaw(t *testing.T) {
	name, err := runtime.Codec.EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = runtime.Codec.EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
