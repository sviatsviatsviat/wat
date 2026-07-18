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
		RawEvent{},
	}
	for _, ev := range cases {
		name := ev.EventName()
		if name == "" {
			name = "RawEvent"
		}
		t.Run(name, func(t *testing.T) {
			_ = EnvelopeOf(ev)
			_ = RawBytes(ev)
		})
	}
}
