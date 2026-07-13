package cursor

import "testing"

func TestEnvelopeAccessorForEveryDecoder(t *testing.T) {
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

	covered := make(map[string]bool, len(cases))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			envelopeAccessorForEvent(tc.ev)
			envelopeAccessorForValue(tc.ptr)
			if name := tc.ev.EventName(); name != "" {
				covered[name] = true
			}
		})
	}

	for name := range decoders {
		if !covered[name] {
			t.Errorf("decoder %q has no envelope accessor test case", name)
		}
	}
}

func TestDecoderRegistryMatchesEventConstants(t *testing.T) {
	want := []string{
		EventSessionStart,
		EventSessionEnd,
		EventBeforeSubmitPrompt,
		EventPreToolUse,
		EventPostToolUse,
		EventPostToolUseFailure,
		EventBeforeShellExecution,
		EventAfterShellExecution,
		EventBeforeMCPExecution,
		EventAfterMCPExecution,
		EventBeforeReadFile,
		EventAfterFileEdit,
		EventSubagentStart,
		EventSubagentStop,
		EventStop,
		EventPreCompact,
		EventAfterAgentResponse,
		EventAfterAgentThought,
		EventBeforeTabFileRead,
		EventAfterTabFileEdit,
		EventWorkspaceOpen,
	}
	if len(decoders) != len(want) {
		t.Fatalf("decoder count = %d, want %d", len(decoders), len(want))
	}
	for _, name := range want {
		if _, ok := decoders[name]; !ok {
			t.Fatalf("missing decoder for %q", name)
		}
	}
}
