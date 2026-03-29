package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestAfterFileEditHookAdapter_viaBuilder_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	adapter, err := NewHookAdapterFromEventFields[AfterFileEditFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestAfterFileEditHookAdapter_viaBuilder_invalidPayload(t *testing.T) {
	mock := cli.NewMockConsole()
	common := HookDataCommon{HookEventName: "afterFileEdit"}
	_, err := NewHookAdapterFromEventFields[AfterFileEditFields](mock, []byte(`not json`), common)
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
	if !strings.Contains(err.Error(), "invalid cursor afterFileEdit payload") {
		t.Fatalf("error should include event-specific prefix: %v", err)
	}
}

func TestAfterFileEditHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterFileEdit","conversation_id":"cid-1","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)
	common, err := NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	a, err := NewHookAdapterFromEventFields[AfterFileEditFields](mock, raw, common)
	if err != nil {
		t.Fatalf("NewHookAdapterFromEventFields: %v", err)
	}
	adapter, ok := a.(*AfterFileEditCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *AfterFileEditCursorHookAdapter, got %T", a)
	}
	if adapter.CommonInput.HookEventName != "afterFileEdit" {
		t.Errorf("HOOK_EVENT_NAME: got %q", adapter.CommonInput.HookEventName)
	}
	if adapter.CommonInput.ConversationID != "cid-1" {
		t.Errorf("CONVERSATION_ID: got %q", adapter.CommonInput.ConversationID)
	}
	if adapter.EventSpecificInput == nil || adapter.EventSpecificInput.FilePath != "D:/repo/file.go" {
		t.Fatalf("EventSpecificInput.FilePath: got %#v", adapter.EventSpecificInput)
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_afterFileEditUsesCursorHookAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"D:/repo/file.go","edits":[{"old_string":"a","new_string":"b"}]}`)

	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	if _, ok := adapter.(*AfterFileEditCursorHookAdapter); !ok {
		t.Fatalf("adapter type: want *AfterFileEditCursorHookAdapter, got %T", adapter)
	}
}

func TestHookAdapterFactory_afterTabFileEditUsesCursorHookAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterTabFileEdit","file_path":"D:/repo/tab.go","edits":[{"old_string":"x","new_string":"y","range":{"start_line_number":1,"start_column":0,"end_line_number":1,"end_column":1},"old_line":"a","new_line":"b"}]}`)

	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	a, ok := adapter.(*AfterFileEditCursorHookAdapter)
	if !ok || a == nil {
		t.Fatalf("adapter type: want *AfterFileEditCursorHookAdapter, got %T", adapter)
	}
	if a.EventSpecificInput == nil || a.EventSpecificInput.FilePath != "D:/repo/tab.go" {
		t.Fatalf("EventSpecificInput.FilePath: got %#v", a.EventSpecificInput)
	}
	if len(a.EventSpecificInput.Edits) != 1 || a.EventSpecificInput.Edits[0].EditRange == nil {
		t.Fatalf("expected one edit with range: %#v", a.EventSpecificInput.Edits)
	}
}
