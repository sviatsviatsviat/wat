package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestAfterAgentResponseHookAdapter_viaBuilder_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterAgentResponse","text":"done"}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter, err := NewHookAdapterFromEventFields[AfterAgentResponseFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestAfterAgentResponseHookAdapter_viaBuilder_invalidPayload(t *testing.T) {
	mock := cli.NewMockConsole()
	common := HookDataCommon{HookEventName: "afterAgentResponse"}
	_, err := NewHookAdapterFromEventFields[AfterAgentResponseFields](mock, []byte(`not json`), common)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor afterAgentResponse payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestAfterAgentThoughtHookAdapter_viaBuilder_invalidPayload(t *testing.T) {
	mock := cli.NewMockConsole()
	common := HookDataCommon{HookEventName: "afterAgentThought"}
	_, err := NewHookAdapterFromEventFields[AfterAgentThoughtFields](mock, []byte(`not json`), common)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor afterAgentThought payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestAfterAgentThoughtHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterAgentThought","conversation_id":"cid-1","text":"think","duration_ms":42}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a, err := NewHookAdapterFromEventFields[AfterAgentThoughtFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	adapter, ok := a.(*AfterAgentThoughtCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *AfterAgentThoughtCursorHookAdapter, got %T", a)
	}
	if adapter.EventSpecificInput == nil || adapter.EventSpecificInput.Text != "think" || adapter.EventSpecificInput.DurationMs != 42 {
		t.Fatalf("EventSpecificInput: got %#v", adapter.EventSpecificInput)
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_afterAgentResponseUsesTypedAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterAgentResponse","text":"hi"}`)
	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	a, ok := adapter.(*AfterAgentResponseCursorHookAdapter)
	if !ok {
		t.Fatalf("adapter type: want *AfterAgentResponseCursorHookAdapter, got %T", adapter)
	}
	if a.EventSpecificInput == nil || a.EventSpecificInput.Text != "hi" {
		t.Fatalf("EventSpecificInput.Text: got %#v", a.EventSpecificInput)
	}
}

func TestHookAdapterFactory_afterAgentThoughtUsesTypedAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterAgentThought","text":"t","duration_ms":99}`)
	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	a, ok := adapter.(*AfterAgentThoughtCursorHookAdapter)
	if !ok {
		t.Fatalf("adapter type: want *AfterAgentThoughtCursorHookAdapter, got %T", adapter)
	}
	if a.EventSpecificInput == nil || a.EventSpecificInput.DurationMs != 99 {
		t.Fatalf("EventSpecificInput.DurationMs: got %#v", a.EventSpecificInput)
	}
}
