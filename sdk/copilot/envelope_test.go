package copilot

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEventNames(t *testing.T) {
	cases := []struct {
		ev       run.Event
		wantName string
	}{
		{SessionStart{}, EventSessionStart},
		{SessionEnd{}, EventSessionEnd},
		{UserPromptSubmitted{}, EventUserPromptSubmitted},
		{PreToolUse{}, EventPreToolUse},
		{PostToolUse{}, EventPostToolUse},
		{PostToolUseFailure{}, EventPostToolUseFailure},
		{PermissionRequest{}, EventPermissionRequest},
		{SubagentStart{}, EventSubagentStart},
		{SubagentStop{}, EventSubagentStop},
		{AgentStop{}, EventAgentStop},
		{PreCompact{}, EventPreCompact},
		{Notification{}, EventNotification},
		{ErrorOccurred{}, EventErrorOccurred},
	}
	for _, tc := range cases {
		t.Run(tc.wantName, func(t *testing.T) {
			if got := tc.ev.EventName(); got != tc.wantName {
				t.Fatalf("EventName() = %q, want %q", got, tc.wantName)
			}
		})
	}
}
