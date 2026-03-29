package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestNewAfterMCPExecutionHookAdapter_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterMCPExecution","tool_name":"t","tool_input":"{}","result_json":"{}","duration":1}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter, err := NewHookAdapterFromEventFields[AfterMCPExecutionFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestAfterMCPExecutionHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterMCPExecution","conversation_id":"cid-1","tool_name":"search","tool_input":"{\"a\":1}","result_json":"{\"b\":2}","duration":1234}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a, err := NewHookAdapterFromEventFields[AfterMCPExecutionFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	adapter, ok := a.(*AfterMCPExecutionCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *AfterMCPExecutionCursorHookAdapter, got %T", a)
	}
	if adapter.CommonInput.HookEventName != "afterMCPExecution" {
		t.Errorf("HOOK_EVENT_NAME: got %q", adapter.CommonInput.HookEventName)
	}
	if adapter.CommonInput.ConversationID != "cid-1" {
		t.Errorf("CONVERSATION_ID: got %q", adapter.CommonInput.ConversationID)
	}
	if adapter.EventSpecificInput == nil {
		t.Fatal("EventSpecificInput must be set")
	}
	if adapter.EventSpecificInput.ToolName != "search" {
		t.Errorf("TOOL_NAME: got %q", adapter.EventSpecificInput.ToolName)
	}
	if adapter.EventSpecificInput.ToolInput != `{"a":1}` {
		t.Errorf("TOOL_INPUT: got %q", adapter.EventSpecificInput.ToolInput)
	}
	if adapter.EventSpecificInput.ResultJSON != `{"b":2}` {
		t.Errorf("RESULT_JSON: got %q", adapter.EventSpecificInput.ResultJSON)
	}
	if adapter.EventSpecificInput.Duration != 1234 {
		t.Errorf("DURATION: got %v", adapter.EventSpecificInput.Duration)
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_afterMCPExecutionUsesCursorHookAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterMCPExecution","tool_name":"x","tool_input":"","result_json":"","duration":0}`)

	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	if _, ok := adapter.(*AfterMCPExecutionCursorHookAdapter); !ok {
		t.Fatalf("adapter type: want *AfterMCPExecutionCursorHookAdapter, got %T", adapter)
	}
}
