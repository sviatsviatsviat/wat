package cursor

import "testing"

func TestEventEnvelopeAccess(t *testing.T) {
	cases := []Event{
		SessionStart{},
		SessionEnd{},
		BeforeSubmitPrompt{},
		PreToolUse{},
		PostToolUse{},
		PostToolUseFailure{},
		BeforeShellExecution{},
		AfterShellExecution{},
		BeforeMCPExecution{},
		AfterMCPExecution{},
		BeforeReadFile{},
		AfterFileEdit{},
		SubagentStart{},
		SubagentStop{},
		Stop{},
		PreCompact{},
		AfterAgentResponse{},
		AfterAgentThought{},
		BeforeTabFileRead{},
		AfterTabFileEdit{},
		WorkspaceOpen{},
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
