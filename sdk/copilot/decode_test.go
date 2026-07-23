package copilot

import (
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte(`{"session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"FutureEvent","session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"PreToolUse","session_id":"s1","cwd":"/w","tool_name":123}`)
	_, err := runtime.Codec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestEventNames(t *testing.T) {
	cases := []struct {
		ev       hookkit.Event
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

func TestEventNameFromRaw(t *testing.T) {
	name, err := runtime.Codec.EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = runtime.Codec.EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
