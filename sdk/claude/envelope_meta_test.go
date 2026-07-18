package claude

import "testing"

func TestEnvelopeAccessorForEveryEvent(t *testing.T) {
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			envelopeAccessorForEvent(tc.ev)
			envelopeAccessorForValue(tc.ptr)
		})
	}
}
