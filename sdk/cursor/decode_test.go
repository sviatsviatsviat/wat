package cursor

import (
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestDecode_RequiresHookEventName(t *testing.T) {
	raw := `{"conversation_id":"c1","command":"ls","cwd":"/w"}`
	_, err := runtime.Codec.Decode([]byte(raw))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
	if !strings.Contains(err.Error(), "hook_event_name") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := runtime.Codec.Decode([]byte(`{"hook_event_name":"FutureEvent","conversation_id":"c1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	raw := []byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":123}`)
	_, err := runtime.Codec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestEventNames(t *testing.T) {
	cases := []run.Event{
		SessionStart{},
		SessionEnd{},
		BeforeSubmitPrompt{},
		PreToolUse{},
		PostToolUse{},
		PostToolUseFailure{},
		BeforeShellExecution{},
		AfterShellExecution{},
		BeforeMCPExecution{},
		AfterMCPExecution{},
		BeforeReadFile{},
		AfterFileEdit{},
		SubagentStart{},
		SubagentStop{},
		Stop{},
		PreCompact{},
		AfterAgentResponse{},
		AfterAgentThought{},
		BeforeTabFileRead{},
		AfterTabFileEdit{},
		WorkspaceOpen{},
	}
	for _, ev := range cases {
		t.Run(ev.EventName(), func(t *testing.T) {
			if ev.EventName() == "" {
				t.Fatal("EventName() empty")
			}
		})
	}
}
