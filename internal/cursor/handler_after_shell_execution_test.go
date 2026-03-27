package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

func TestNewAfterShellExecutionHookHandler_success(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"npm test","output":"ok","duration":1,"sandbox":false}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterShellExecutionFields](raw, common)
	if err != nil {
		t.Fatalf("NewHookHandlerFromEventFields: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestAfterShellExecutionHookHandler_Handle_wiresContextAndOutput(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterShellExecution","conversation_id":"cid-1","command":"npm test","output":"all good","duration":1234,"sandbox":true}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterShellExecutionFields](raw, common)
	if err != nil {
		t.Fatalf("NewHookHandlerFromEventFields: %v", err)
	}

	var seenCtx *core.HookContext
	hookCommand := stubHookCommand{execute: func(ctx *core.HookContext) int {
		seenCtx = ctx
		if ctx.HookHost != cursorcore.HookHostCursor {
			t.Errorf("HookHost: want %q, got %q", cursorcore.HookHostCursor, ctx.HookHost)
		}
		rd, ok := ctx.ParsedData.(*cursorcore.CursorHookRunData[cursorcore.AfterShellExecutionFields])
		if !ok || rd == nil {
			t.Fatalf("ParsedData must be *CursorHookRunData[AfterShellExecutionFields], got %T", ctx.ParsedData)
		}
		if rd.Common.HookEventName != "afterShellExecution" {
			t.Errorf("HOOK_EVENT_NAME: got %q", rd.Common.HookEventName)
		}
		if rd.Common.ConversationID != "cid-1" {
			t.Errorf("CONVERSATION_ID: got %q", rd.Common.ConversationID)
		}
		if rd.EventSpecific == nil {
			t.Fatal("EventSpecific must be set")
		}
		if rd.EventSpecific.Command != "npm test" {
			t.Errorf("COMMAND: got %q", rd.EventSpecific.Command)
		}
		if rd.EventSpecific.Output != "all good" {
			t.Errorf("OUTPUT: got %q", rd.EventSpecific.Output)
		}
		if rd.EventSpecific.Duration != 1234 {
			t.Errorf("DURATION: got %v", rd.EventSpecific.Duration)
		}
		if !rd.EventSpecific.Sandbox {
			t.Error("SANDBOX: want true")
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
	if result.Output != cursorcore.DefaultHookResponseLine {
		t.Fatalf("output: want %q, got %q", cursorcore.DefaultHookResponseLine, result.Output)
	}
}

func TestHookHandlerFactory_afterShellExecutionUsesCursorHookHandler(t *testing.T) {
	factory := NewHookHandlerFactory()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"x","output":"","duration":0,"sandbox":false}`)

	handler, err := factory.HookHandlerFromJSON(raw)
	if err != nil {
		t.Fatalf("HookHandlerFromJSON: %v", err)
	}
	if _, ok := handler.(cursorcore.CursorHookHandler[cursorcore.AfterShellExecutionFields]); !ok {
		t.Fatalf("handler type: want CursorHookHandler[AfterShellExecutionFields], got %T", handler)
	}
}
