package cursorcore

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
		if ctx.TemplateBindings == nil {
			t.Error("HookContext.TemplateBindings must be set")
			return 1
		}
		hookEventName, ok := ctx.TemplateBindings.TemplateValue("HOOK_EVENT_NAME")
		if !ok || hookEventName != "afterFileEdit" {
			t.Errorf("HOOK_EVENT_NAME: want afterFileEdit, ok=%v got=%q", ok, hookEventName)
		}
		conversationID, ok := ctx.TemplateBindings.TemplateValue("CONVERSATION_ID")
		if !ok || conversationID != "cid-1" {
			t.Errorf("CONVERSATION_ID: want cid-1, ok=%v got=%q", ok, conversationID)
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
	if result.Output != defaultHookResponseLine {
		t.Fatalf("output: want %q, got %q", defaultHookResponseLine, result.Output)
	}
}

type stubHookCommand struct {
	execute func(*core.HookContext) int
}

func (stub stubHookCommand) Execute(ctx *core.HookContext) int {
	if stub.execute == nil {
		return 0
	}
	return stub.execute(ctx)
}
