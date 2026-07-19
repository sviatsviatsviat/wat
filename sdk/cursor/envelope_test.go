package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEventNames(t *testing.T) {
	cases := []run.Event{
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
			if ev.EventName() == "" {
				t.Fatal("EventName() empty")
			}
		})
	}
}
