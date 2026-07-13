package claude

import "testing"

func TestEnvelopeAccessorForEveryDecoder(t *testing.T) {
	cases := []struct {
		name string
		ev   Event
		ptr  envelopeAccessor
	}{
		{"SessionStart", SessionStart{}, &SessionStart{}},
		{"Setup", Setup{}, &Setup{}},
		{"SessionEnd", SessionEnd{}, &SessionEnd{}},
		{"UserPromptSubmit", UserPromptSubmit{}, &UserPromptSubmit{}},
		{"UserPromptExpansion", UserPromptExpansion{}, &UserPromptExpansion{}},
		{"PreToolUse", PreToolUse{}, &PreToolUse{}},
		{"PostToolUse", PostToolUse{}, &PostToolUse{}},
		{"PostToolUseFailure", PostToolUseFailure{}, &PostToolUseFailure{}},
		{"PostToolBatch", PostToolBatch{}, &PostToolBatch{}},
		{"PermissionRequest", PermissionRequest{}, &PermissionRequest{}},
		{"PermissionDenied", PermissionDenied{}, &PermissionDenied{}},
		{"SubagentStart", SubagentStart{}, &SubagentStart{}},
		{"SubagentStop", SubagentStop{}, &SubagentStop{}},
		{"TaskCreated", TaskCreated{}, &TaskCreated{}},
		{"TaskCompleted", TaskCompleted{}, &TaskCompleted{}},
		{"Stop", Stop{}, &Stop{}},
		{"StopFailure", StopFailure{}, &StopFailure{}},
		{"TeammateIdle", TeammateIdle{}, &TeammateIdle{}},
		{"Notification", Notification{}, &Notification{}},
		{"MessageDisplay", MessageDisplay{}, &MessageDisplay{}},
		{"InstructionsLoaded", InstructionsLoaded{}, &InstructionsLoaded{}},
		{"ConfigChange", ConfigChange{}, &ConfigChange{}},
		{"CwdChanged", CwdChanged{}, &CwdChanged{}},
		{"FileChanged", FileChanged{}, &FileChanged{}},
		{"WorktreeCreate", WorktreeCreate{}, &WorktreeCreate{}},
		{"WorktreeRemove", WorktreeRemove{}, &WorktreeRemove{}},
		{"PreCompact", PreCompact{}, &PreCompact{}},
		{"PostCompact", PostCompact{}, &PostCompact{}},
		{"Elicitation", Elicitation{}, &Elicitation{}},
		{"ElicitationResult", ElicitationResult{}, &ElicitationResult{}},
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
		EventSetup,
		EventSessionEnd,
		EventUserPromptSubmit,
		EventUserPromptExpansion,
		EventPreToolUse,
		EventPostToolUse,
		EventPostToolUseFailure,
		EventPostToolBatch,
		EventPermissionRequest,
		EventPermissionDenied,
		EventSubagentStart,
		EventSubagentStop,
		EventTaskCreated,
		EventTaskCompleted,
		EventStop,
		EventStopFailure,
		EventTeammateIdle,
		EventNotification,
		EventMessageDisplay,
		EventInstructionsLoaded,
		EventConfigChange,
		EventCwdChanged,
		EventFileChanged,
		EventWorktreeCreate,
		EventWorktreeRemove,
		EventPreCompact,
		EventPostCompact,
		EventElicitation,
		EventElicitationResult,
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
