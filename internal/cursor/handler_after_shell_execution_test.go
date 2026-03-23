package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestNewAfterShellExecutionHookHandler_success(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"npm test","output":"ok","duration":1,"sandbox":false}`)
	hookDataCommon, err := newHookDataCommon(raw)
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	handler, err := newAfterShellExecutionHookHandler(raw, hookDataCommon)
	if err != nil {
		t.Fatalf("newAfterShellExecutionHookHandler: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestAfterShellExecutionHookHandler_Handle_wiresContextAndOutput(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterShellExecution","conversation_id":"cid-1","command":"npm test","output":"all good","duration":1234,"sandbox":true}`)
	hookDataCommon, err := newHookDataCommon(raw)
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	handler, err := newAfterShellExecutionHookHandler(raw, hookDataCommon)
	if err != nil {
		t.Fatalf("newAfterShellExecutionHookHandler: %v", err)
	}

	var seenCtx *core.HookContext
	hookCommand := stubHookCommand{execute: func(ctx *core.HookContext) int {
		seenCtx = ctx
		if ctx.TemplateBindings == nil {
			t.Error("HookContext.TemplateBindings must be set")
			return 1
		}
		assertTemplateBindingValue(t, ctx.TemplateBindings, "HOOK_EVENT_NAME", "afterShellExecution")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "CONVERSATION_ID", "cid-1")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "COMMAND", "npm test")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "OUTPUT", "all good")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "DURATION", "1234")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "SANDBOX", "true")
		return 42
	}}

	result := handler.Handle(hookCommand)
	if result.Code != 42 {
		t.Fatalf("Handle exit code: want 42, got %d", result.Code)
	}
	if seenCtx == nil {
		t.Fatal("Command.Execute was not called")
	}
	if result.Output != defaultHookResponseLine {
		t.Fatalf("output: want %q, got %q", defaultHookResponseLine, result.Output)
	}
}

func TestHookHandlerFactory_afterShellExecutionUsesDedicatedHandler(t *testing.T) {
	factory := NewHookHandlerFactory()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"x","output":"","duration":0,"sandbox":false}`)

	handler, err := factory.HookHandlerFromJSON(raw)
	if err != nil {
		t.Fatalf("HookHandlerFromJSON: %v", err)
	}
	if _, ok := handler.(afterShellExecutionHookHandler); !ok {
		t.Fatalf("handler type: want afterShellExecutionHookHandler, got %T", handler)
	}
}
