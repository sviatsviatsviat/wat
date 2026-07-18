package cursor

import "testing"

func TestEnvelopeAccessorForEveryEvent(t *testing.T) {
	cases := []struct {
		name string
		ev   Event
		ptr  envelopeAccessor
	}{
		{"SessionStart", SessionStart{}, &SessionStart{}},
		{"SessionEnd", SessionEnd{}, &SessionEnd{}},
		{"BeforeSubmitPrompt", BeforeSubmitPrompt{}, &BeforeSubmitPrompt{}},
		{"PreToolUse", PreToolUse{}, &PreToolUse{}},
		{"PostToolUse", PostToolUse{}, &PostToolUse{}},
		{"PostToolUseFailure", PostToolUseFailure{}, &PostToolUseFailure{}},
		{"BeforeShellExecution", BeforeShellExecution{}, &BeforeShellExecution{}},
		{"AfterShellExecution", AfterShellExecution{}, &AfterShellExecution{}},
		{"BeforeMCPExecution", BeforeMCPExecution{}, &BeforeMCPExecution{}},
		{"AfterMCPExecution", AfterMCPExecution{}, &AfterMCPExecution{}},
		{"BeforeReadFile", BeforeReadFile{}, &BeforeReadFile{}},
		{"AfterFileEdit", AfterFileEdit{}, &AfterFileEdit{}},
		{"SubagentStart", SubagentStart{}, &SubagentStart{}},
		{"SubagentStop", SubagentStop{}, &SubagentStop{}},
		{"Stop", Stop{}, &Stop{}},
		{"PreCompact", PreCompact{}, &PreCompact{}},
		{"AfterAgentResponse", AfterAgentResponse{}, &AfterAgentResponse{}},
		{"AfterAgentThought", AfterAgentThought{}, &AfterAgentThought{}},
		{"BeforeTabFileRead", BeforeTabFileRead{}, &BeforeTabFileRead{}},
		{"AfterTabFileEdit", AfterTabFileEdit{}, &AfterTabFileEdit{}},
		{"WorkspaceOpen", WorkspaceOpen{}, &WorkspaceOpen{}},
		{"RawEvent", RawEvent{}, &RawEvent{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			envelopeAccessorForEvent(tc.ev)
			envelopeAccessorForValue(tc.ptr)
		})
	}
}
