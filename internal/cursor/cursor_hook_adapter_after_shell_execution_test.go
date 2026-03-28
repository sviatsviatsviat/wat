package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestNewAfterShellExecutionHookAdapter_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"npm test","output":"ok","duration":1,"sandbox":false}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter, err := NewHookAdapterFromEventFields[AfterShellExecutionFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestAfterShellExecutionHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterShellExecution","conversation_id":"cid-1","command":"npm test","output":"all good","duration":1234,"sandbox":true}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a, err := NewHookAdapterFromEventFields[AfterShellExecutionFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	adapter, ok := a.(*AfterShellExecutionCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *AfterShellExecutionCursorHookAdapter, got %T", a)
	}
	if adapter.CommonInput.HookEventName != "afterShellExecution" {
		t.Errorf("HOOK_EVENT_NAME: got %q", adapter.CommonInput.HookEventName)
	}
	if adapter.CommonInput.ConversationID != "cid-1" {
		t.Errorf("CONVERSATION_ID: got %q", adapter.CommonInput.ConversationID)
	}
	if adapter.EventSpecificInput == nil {
		t.Fatal("EventSpecificInput must be set")
	}
	if adapter.EventSpecificInput.Command != "npm test" {
		t.Errorf("COMMAND: got %q", adapter.EventSpecificInput.Command)
	}
	if adapter.EventSpecificInput.Output != "all good" {
		t.Errorf("OUTPUT: got %q", adapter.EventSpecificInput.Output)
	}
	if adapter.EventSpecificInput.Duration != 1234 {
		t.Errorf("DURATION: got %v", adapter.EventSpecificInput.Duration)
	}
	if !adapter.EventSpecificInput.Sandbox {
		t.Error("SANDBOX: want true")
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_afterShellExecutionUsesCursorHookAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"x","output":"","duration":0,"sandbox":false}`)

	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	if _, ok := adapter.(*AfterShellExecutionCursorHookAdapter); !ok {
		t.Fatalf("adapter type: want *AfterShellExecutionCursorHookAdapter, got %T", adapter)
	}
}
