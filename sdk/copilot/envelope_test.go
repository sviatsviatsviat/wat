package copilot

import "testing"

func TestEventEnvelopeAccess(t *testing.T) {
	cases := []Event{
		SessionStart{},
		SessionEnd{},
		UserPromptSubmitted{},
		PreToolUse{},
		PostToolUse{},
		PostToolUseFailure{},
		PermissionRequest{},
		SubagentStart{},
		SubagentStop{},
		AgentStop{},
		PreCompact{},
		Notification{},
		ErrorOccurred{},
	}
	for _, ev := range cases {
		t.Run(ev.EventName(), func(t *testing.T) {
			_ = EnvelopeOf(ev)
		})
	}
}
