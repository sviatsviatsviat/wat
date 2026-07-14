package copilot

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal"
)

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

	for _, name := range internal.RegisteredDecoders() {
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
	if len(internal.RegisteredDecoders()) != len(want) {
		t.Fatalf("decoder count = %d, want %d", len(internal.RegisteredDecoders()), len(want))
	}
	for _, name := range want {
		if _, ok := internal.DecoderFor(name); !ok {
			t.Fatalf("missing decoder for %q", name)
		}
	}
}
