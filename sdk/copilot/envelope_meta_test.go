package copilot

import "testing"

func TestEnvelopeAccessorForEveryEvent(t *testing.T) {
	cases := []struct {
		name string
		ev   Event
		ptr  envelopeAccessor
	}{
		{"SessionStart", SessionStart{}, &SessionStart{}},
		{"SessionEnd", SessionEnd{}, &SessionEnd{}},
		{"UserPromptSubmitted", UserPromptSubmitted{}, &UserPromptSubmitted{}},
		{"PreToolUse", PreToolUse{}, &PreToolUse{}},
		{"PostToolUse", PostToolUse{}, &PostToolUse{}},
		{"PostToolUseFailure", PostToolUseFailure{}, &PostToolUseFailure{}},
		{"PermissionRequest", PermissionRequest{}, &PermissionRequest{}},
		{"SubagentStart", SubagentStart{}, &SubagentStart{}},
		{"SubagentStop", SubagentStop{}, &SubagentStop{}},
		{"AgentStop", AgentStop{}, &AgentStop{}},
		{"PreCompact", PreCompact{}, &PreCompact{}},
		{"Notification", Notification{}, &Notification{}},
		{"ErrorOccurred", ErrorOccurred{}, &ErrorOccurred{}},
		{"RawEvent", RawEvent{}, &RawEvent{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			envelopeAccessorForEvent(tc.ev)
			envelopeAccessorForValue(tc.ptr)
		})
	}
}
