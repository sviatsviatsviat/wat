package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestNewDefaultHookAdapter_success(t *testing.T) {
	mock := cli.NewMockConsole()
	hookData, err := NewHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit"}`))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter := NewDefaultHookAdapter(mock, hookData)
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestDefaultHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	hookData, err := NewHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit","conversation_id":"cid-1"}`))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a := NewDefaultHookAdapter(mock, hookData)
	adapter, ok := a.(*DefaultCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *DefaultCursorHookAdapter, got %T", a)
	}
	if adapter.HookHost() != HookHostCursor {
		t.Errorf("HookHost: want %q, got %q", HookHostCursor, adapter.HookHost())
	}
	if adapter.CommonInput.HookEventName != "afterFileEdit" {
		t.Errorf("HOOK_EVENT_NAME: want afterFileEdit, got %q", adapter.CommonInput.HookEventName)
	}
	if adapter.CommonInput.ConversationID != "cid-1" {
		t.Errorf("CONVERSATION_ID: want cid-1, got %q", adapter.CommonInput.ConversationID)
	}
	if adapter.EventSpecificInput != nil {
		t.Error("default adapter must leave EventSpecificInput nil")
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}
