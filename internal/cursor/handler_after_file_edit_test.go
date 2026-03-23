package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor/core"
)

func TestNewCursorEventHookHandlerBuilder_afterFileEdit_success(t *testing.T) {
	build := cursorcore.NewCursorEventHookHandlerBuilder(afterFileEditPlaceholderExtractors)
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := build(raw, common)
	if err != nil {
		t.Fatalf("hookHandlerBuilder: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestNewCursorEventHookHandlerBuilder_afterFileEdit_invalidPayload(t *testing.T) {
	build := cursorcore.NewCursorEventHookHandlerBuilder(afterFileEditPlaceholderExtractors)
	common := cursorcore.HookDataCommon{HookEventName: "afterFileEdit"}
	_, err := build([]byte(`not json`), common)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor afterFileEdit payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestAfterFileEditHookHandler_Handle_wiresContextAndOutput(t *testing.T) {
	build := cursorcore.NewCursorEventHookHandlerBuilder(afterFileEditPlaceholderExtractors)
	raw := []byte(`{"hook_event_name":"afterFileEdit","conversation_id":"cid-1","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := build(raw, common)
	if err != nil {
		t.Fatalf("hookHandlerBuilder: %v", err)
	}

	var seenCtx *core.HookContext
	hookCommand := stubHookCommand{execute: func(ctx *core.HookContext) int {
		seenCtx = ctx
		if ctx.TemplateBindings == nil {
			t.Error("HookContext.TemplateBindings must be set")
			return 1
		}
		assertTemplateBindingValue(t, ctx.TemplateBindings, "HOOK_EVENT_NAME", "afterFileEdit")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "CONVERSATION_ID", "cid-1")
		assertTemplateBindingValue(t, ctx.TemplateBindings, "FILE_PATH", "D:/repo/file.go")
		return 42
	}}

	result := handler.Handle(hookCommand)
	if result.Code != 42 {
		t.Fatalf("Handle exit code: want 42, got %d", result.Code)
	}
	if seenCtx == nil {
		t.Fatal("Command.Execute was not called")
	}
	if result.Output != "{}\n" {
		t.Fatalf("output: want %q, got %q", "{}\n", result.Output)
	}
}

func TestHookHandlerFactory_afterFileEditUsesDedicatedHandler(t *testing.T) {
	factory := NewHookHandlerFactory()
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)

	handler, err := factory.HookHandlerFromJSON(raw)
	if err != nil {
		t.Fatalf("HookHandlerFromJSON: %v", err)
	}
	if _, ok := handler.(cursorcore.EventHookHandler[cursorcore.HookDataWithCommon[hookDataAfterFileEditFields]]); !ok {
		t.Fatalf("handler type: want EventHookHandler for afterFileEdit fields, got %T", handler)
	}
}
