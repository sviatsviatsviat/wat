package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cursor/core"
)

func TestAfterFileEditPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"FILE_PATH": {},
	}
	if len(afterFileEditPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(afterFileEditPlaceholderExtractors))
	}
	for placeholderKey := range wantKeys {
		if _, ok := afterFileEditPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsAfterFileEdit_templateValueEventAndCommonFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "conv-1",
		"file_path": "D:/repo/file.go",
		"edits": [{"old_string":"a","new_string":"b"}]
	}`
	commonData, err := cursorcore.NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := cursorcore.NewHookDataWithCommon[hookDataAfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	bindings := cursorcore.NewTemplateBindingsEvent(hookData.HookDataCommon, hookData.Fields, afterFileEditPlaceholderExtractors)

	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterFileEdit")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "FILE_PATH", "D:/repo/file.go")
}

func TestTemplateBindingsAfterFileEdit_editsPlaceholderNotDefined(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"x","edits":[{"old_string":"a","new_string":"b"}]}`)
	commonData, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := cursorcore.NewHookDataWithCommon[hookDataAfterFileEditFields](raw, commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	bindings := cursorcore.NewTemplateBindingsEvent(hookData.HookDataCommon, hookData.Fields, afterFileEditPlaceholderExtractors)

	_, ok := bindings.TemplateValue("EDITS")
	if ok {
		t.Fatal("EDITS must not be a defined placeholder")
	}
	_, ok = bindings.TemplateValue("OLD_STRING")
	if ok {
		t.Fatal("OLD_STRING must not be a defined placeholder")
	}
	_, ok = bindings.TemplateValue("NEW_STRING")
	if ok {
		t.Fatal("NEW_STRING must not be a defined placeholder")
	}
}

func TestTemplateBindingsAfterFileEdit_unknownKey(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterFileEdit","file_path":"x"}`)
	commonData, err := cursorcore.NewHookDataCommon(raw)
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := cursorcore.NewHookDataWithCommon[hookDataAfterFileEditFields](raw, commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	bindings := cursorcore.NewTemplateBindingsEvent(hookData.HookDataCommon, hookData.Fields, afterFileEditPlaceholderExtractors)

	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}
