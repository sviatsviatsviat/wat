package claude

import (
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
