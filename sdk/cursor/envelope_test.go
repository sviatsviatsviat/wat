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
	}
	for _, ev := range cases {
		t.Run(ev.EventName(), func(t *testing.T) {
			_ = EnvelopeOf(ev)
		})
	}
}
