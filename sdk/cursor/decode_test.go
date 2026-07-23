package cursor

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/pretooluse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_RequiresHookEventName(t *testing.T) {
	raw := `{"conversation_id":"c1","command":"ls","cwd":"/w"}`
	_, err := testCodec.Decode([]byte(raw))
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
	_, err := testCodec.Decode([]byte(`{"hook_event_name":"FutureEvent","conversation_id":"c1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_InvalidTypedEvent(t *testing.T) {
	pretooluse.RegisterHandler(hookkit.NewDialect(testCodec), func(_ context.Context, _ pretooluse.Event, r event.PermissionResults) (event.PermissionOutput, error) {
		return r.Noop(), nil
	})
	raw := []byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":123}`)
	_, err := testCodec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestEventNames(t *testing.T) {
	cases := []hookkit.Event{
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
