package claude

import "testing"

func TestEventEnvelopeAccess(t *testing.T) {
	cases := []Event{
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
		RawEvent{},
	}
	for _, ev := range cases {
		t.Run(ev.EventName(), func(t *testing.T) {
			_ = EnvelopeOf(ev)
		})
	}
}
