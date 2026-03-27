package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

func TestAfterFileEditHookHandler_viaBuilder_success(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	handler, err := cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterFileEditFields](raw, common)
	if err != nil {
		t.Fatalf("NewHookHandlerFromEventFields: %v", err)
	}
	if handler == nil {
		t.Fatal("expected non-nil HookHandler")
	}
}

func TestAfterFileEditHookHandler_viaBuilder_invalidPayload(t *testing.T) {
	common := cursorcore.HookDataCommon{HookEventName: "afterFileEdit"}
	_, err := cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterFileEditFields]([]byte(`not json`), common)
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
	handler, err := cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterFileEditFields](raw, common)
	if err != nil {
		t.Fatalf("NewHookHandlerFromEventFields: %v", err)
	}

	var seenCtx *core.HookContext
	hookCommand := stubHookCommand{execute: func(ctx *core.HookContext) int {
		seenCtx = ctx
		if ctx.HookHost != cursorcore.HookHostCursor {
			t.Errorf("HookHost: want %q, got %q", cursorcore.HookHostCursor, ctx.HookHost)
		}
		rd, ok := ctx.ParsedData.(*cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields])
		if !ok || rd == nil {
			t.Fatalf("ParsedData must be *CursorHookRunData[AfterFileEditFields], got %T", ctx.ParsedData)
		}
		if rd.Common.HookEventName != "afterFileEdit" {
			t.Errorf("HOOK_EVENT_NAME: got %q", rd.Common.HookEventName)
		}
		if rd.Common.ConversationID != "cid-1" {
			t.Errorf("CONVERSATION_ID: got %q", rd.Common.ConversationID)
		}
		if rd.EventSpecific == nil || rd.EventSpecific.FilePath != "D:/repo/file.go" {
			t.Fatalf("EventSpecific.FilePath: got %#v", rd.EventSpecific)
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

func TestHookHandlerFactory_afterFileEditUsesCursorHookHandler(t *testing.T) {
	factory := NewHookHandlerFactory()
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)

	handler, err := factory.HookHandlerFromJSON(raw)
	if err != nil {
		t.Fatalf("HookHandlerFromJSON: %v", err)
	}
	if _, ok := handler.(cursorcore.CursorHookHandler[cursorcore.AfterFileEditFields]); !ok {
		t.Fatalf("handler type: want CursorHookHandler[AfterFileEditFields], got %T", handler)
	}
}
