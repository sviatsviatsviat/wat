package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor/core"
)

func testWatExecCtx() core.WatExecutionContext {
	return core.NewWatExecutionContext("cursor").WithSubcommand("run")
}

func TestNewAfterFileEditHookHandler_success(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := newAfterFileEditHookHandler(raw, common, testWatExecCtx())
	if err != nil {
		t.Fatalf("newAfterFileEditHookHandler: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestNewAfterFileEditHookHandler_invalidPayload(t *testing.T) {
	common := cursorcore.HookDataCommon{HookEventName: "afterFileEdit"}
	_, err := newAfterFileEditHookHandler([]byte(`not json`), common, testWatExecCtx())
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor afterFileEdit payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestAfterFileEditHookHandler_Handle_wiresContextAndOutput(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","conversation_id":"cid-1","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := newAfterFileEditHookHandler(raw, common, testWatExecCtx())
	if err != nil {
		t.Fatalf("newAfterFileEditHookHandler: %v", err)
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
	if result.Output != cursorcore.DefaultHookResponseLine {
		t.Fatalf("output: want %q, got %q", cursorcore.DefaultHookResponseLine, result.Output)
	}
}

func TestHookHandlerFactory_afterFileEditUsesDedicatedHandler(t *testing.T) {
	factory := NewHookHandlerFactory(testWatExecCtx())
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)

	handler, err := factory.HookHandlerFromJSON(raw)
	if err != nil {
		t.Fatalf("HookHandlerFromJSON: %v", err)
	}
	if _, ok := handler.(afterFileEditHookHandler); !ok {
		t.Fatalf("handler type: want afterFileEditHookHandler, got %T", handler)
	}
}

func TestAfterFileEditHookHandler_Handle_filePatternNoMatchSkipsCommand(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.txt","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := newAfterFileEditHookHandler(raw, common, testWatExecCtx().WithFilePattern(`[.]go$`))
	if err != nil {
		t.Fatalf("newAfterFileEditHookHandler: %v", err)
	}

	executed := false
	cmd := stubHookCommand{execute: func(_ *core.HookContext) int {
		executed = true
		return 99
	}}
	result := handler.Handle(cmd)
	if executed {
		t.Fatal("Command.Execute should be skipped when file path does not match regexp")
	}
	if result.Code != 0 {
		t.Fatalf("result code: want 0, got %d", result.Code)
	}
	if result.Output != cursorcore.DefaultHookResponseLine {
		t.Fatalf("output: want %q, got %q", cursorcore.DefaultHookResponseLine, result.Output)
	}
}

func TestAfterFileEditHookHandler_Handle_filePatternMatchExecutesCommand(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := newAfterFileEditHookHandler(raw, common, testWatExecCtx().WithFilePattern(`[.]go$`))
	if err != nil {
		t.Fatalf("newAfterFileEditHookHandler: %v", err)
	}

	executed := false
	cmd := stubHookCommand{execute: func(ctx *core.HookContext) int {
		executed = true
		if ctx.TemplateBindings == nil {
			t.Fatal("HookContext.TemplateBindings must be set")
		}
		return 7
	}}
	result := handler.Handle(cmd)
	if !executed {
		t.Fatal("Command.Execute should run when file path matches regexp")
	}
	if result.Code != 7 {
		t.Fatalf("result code: want 7, got %d", result.Code)
	}
}

func TestNewAfterFileEditHookHandler_invalidFilePatternRegexp(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	_, err = newAfterFileEditHookHandler(raw, common, testWatExecCtx().WithFilePattern(`(`))
	if err == nil {
		t.Fatal("expected error for invalid --file-pattern regexp")
	}
	if !strings.Contains(err.Error(), "invalid --file-pattern regexp") {
		t.Fatalf("error should mention invalid file-pattern regexp: %v", err)
	}
}
