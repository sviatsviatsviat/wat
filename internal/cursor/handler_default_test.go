package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestNewDefaultHookHandler_success(t *testing.T) {
	hookData, err := NewHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit"}`))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := NewDefaultHookHandler(hookData)
	if err != nil {
		t.Fatalf("NewDefaultHookHandler: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestDefaultHookHandler_Handle_wiresContextAndOutput(t *testing.T) {
	hookData, err := NewHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit","conversation_id":"cid-1"}`))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := NewDefaultHookHandler(hookData)
	if err != nil {
		t.Fatal(err)
	}

	var seenCtx *core.HookContext
	hookCommand := stubHookCommand{execute: func(ctx *core.HookContext) int {
		seenCtx = ctx
		if ctx.HookHost != HookHostCursor {
			t.Errorf("HookHost: want %q, got %q", HookHostCursor, ctx.HookHost)
		}
		rd, ok := ctx.ParsedData.(*CursorHookRunData[struct{}])
		if !ok || rd == nil {
			t.Fatalf("ParsedData must be *CursorHookRunData[struct{}], got %T", ctx.ParsedData)
		}
		if rd.Common.HookEventName != "afterFileEdit" {
			t.Errorf("HOOK_EVENT_NAME: want afterFileEdit, got %q", rd.Common.HookEventName)
		}
		if rd.Common.ConversationID != "cid-1" {
			t.Errorf("CONVERSATION_ID: want cid-1, got %q", rd.Common.ConversationID)
		}
		if rd.EventSpecific != nil {
			t.Error("default handler must leave EventSpecific nil")
		}
		return 42
	}}

	result := handler.Handle(hookCommand)
	if result.Code != 42 {
		t.Fatalf("Handle exit code: want 42, got %d", result.Code)
	}
	if seenCtx == nil {
		t.Fatal("Command.Execute was not called")
	}
	if result.Output != cursorHookStdoutSuccessLine {
		t.Fatalf("output: want %q, got %q", cursorHookStdoutSuccessLine, result.Output)
	}
}
