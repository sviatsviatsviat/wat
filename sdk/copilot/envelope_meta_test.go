package copilot

import "testing"

func TestEnvelopeAccessorForEveryDecoder(t *testing.T) {
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
		EventUserPromptSubmitted,
		EventPreToolUse,
		EventPostToolUse,
		EventPostToolUseFailure,
		EventPermissionRequest,
		EventSubagentStart,
		EventSubagentStop,
		EventAgentStop,
		EventPreCompact,
		EventNotification,
		EventErrorOccurred,
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
