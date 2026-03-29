package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestSessionEndHookAdapter_viaBuilder_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"sessionEnd","session_id":"s1"}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter, err := NewHookAdapterFromEventFields[SessionEndFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestSessionEndHookAdapter_viaBuilder_invalidPayload(t *testing.T) {
	mock := cli.NewMockConsole()
	common := HookDataCommon{HookEventName: "sessionEnd"}
	_, err := NewHookAdapterFromEventFields[SessionEndFields](mock, []byte(`not json`), common)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor sessionEnd payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestSessionEndHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{
		"hook_event_name":"sessionEnd",
		"conversation_id":"cid-1",
		"session_id":"sess-99",
		"reason":"completed",
		"duration_ms":45000,
		"is_background_agent":true,
		"final_status":"ok",
		"error_message":""
	}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a, err := NewHookAdapterFromEventFields[SessionEndFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	adapter, ok := a.(*SessionEndCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *SessionEndCursorHookAdapter, got %T", a)
	}
	ev := adapter.EventSpecificInput
	if ev == nil {
		t.Fatal("EventSpecificInput: want non-nil")
	}
	if ev.SessionID != "sess-99" || ev.Reason != "completed" || ev.DurationMs != 45000 ||
		!ev.IsBackgroundAgent || ev.FinalStatus != "ok" {
		t.Fatalf("EventSpecificInput: got %#v", ev)
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_sessionEndUsesTypedAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"sessionEnd","session_id":"s-min"}`)
	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	a, ok := adapter.(*SessionEndCursorHookAdapter)
	if !ok {
		t.Fatalf("adapter type: want *SessionEndCursorHookAdapter, got %T", adapter)
	}
	if a.EventSpecificInput == nil || a.EventSpecificInput.SessionID != "s-min" {
		t.Fatalf("EventSpecificInput.SessionID: got %#v", a.EventSpecificInput)
	}
}

func TestNewHookDataWithCommon_sessionEnd_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "sessionEnd",
		"conversation_id": "c1",
		"session_id": "sid-1",
		"reason": "error",
		"duration_ms": 100,
		"is_background_agent": false,
		"final_status": "failed",
		"error_message": "boom"
	}`
	common, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	d, err := NewHookDataWithCommon[SessionEndFields]([]byte(input), common)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	if d.Fields.SessionID != "sid-1" || d.Fields.Reason != "error" || d.Fields.DurationMs != 100 ||
		d.Fields.IsBackgroundAgent || d.Fields.FinalStatus != "failed" || d.Fields.ErrorMessage != "boom" {
		t.Fatalf("Fields: got %#v", d.Fields)
	}
}
