package claude

import (
	"errors"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEventNames(t *testing.T) {
	cases := []run.Event{
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
	name, err := codec.EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = codec.EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
