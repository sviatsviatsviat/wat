package copilot

import "testing"

func TestEventEnvelopeAccess(t *testing.T) {
	cases := []struct {
		ev       Event
		wantName string
	}{
		{SessionStart{Envelope: Envelope{HookEventName: EventSessionStart}}, EventSessionStart},
		{SessionEnd{Envelope: Envelope{HookEventName: EventSessionEnd}}, EventSessionEnd},
		{UserPromptSubmitted{Envelope: Envelope{HookEventName: EventUserPromptSubmitted}}, EventUserPromptSubmitted},
		{PreToolUse{Envelope: Envelope{HookEventName: EventPreToolUse}}, EventPreToolUse},
		{PostToolUse{Envelope: Envelope{HookEventName: EventPostToolUse}}, EventPostToolUse},
		{PostToolUseFailure{Envelope: Envelope{HookEventName: EventPostToolUseFailure}}, EventPostToolUseFailure},
		{PermissionRequest{Envelope: Envelope{HookEventName: EventPermissionRequest}}, EventPermissionRequest},
		{SubagentStart{Envelope: Envelope{HookEventName: EventSubagentStart}}, EventSubagentStart},
		{SubagentStop{Envelope: Envelope{HookEventName: EventSubagentStop}}, EventSubagentStop},
		{AgentStop{Envelope: Envelope{HookEventName: EventAgentStop}}, EventAgentStop},
		{PreCompact{Envelope: Envelope{HookEventName: EventPreCompact}}, EventPreCompact},
		{Notification{Envelope: Envelope{HookEventName: EventNotification}}, EventNotification},
		{ErrorOccurred{Envelope: Envelope{HookEventName: EventErrorOccurred}}, EventErrorOccurred},
	}
	for _, tc := range cases {
		t.Run(tc.ev.EventName(), func(t *testing.T) {
			env := EnvelopeOf(tc.ev)
			if env.HookEventName != tc.ev.EventName() {
				t.Fatalf("HookEventName = %q, want %q", env.HookEventName, tc.ev.EventName())
			}
			if env.HookEventName != tc.wantName {
				t.Fatalf("HookEventName = %q, want %q", env.HookEventName, tc.wantName)
			}
		})
	}
}
