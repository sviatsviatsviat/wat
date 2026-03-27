package run

import (
	"testing"

	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
	"github.com/sviatsviatsviat/wat/internal/template"
)

func TestTemplateBindingsCommon_templateValueAllCommonFields(t *testing.T) {
	hookData := cursorcore.HookDataCommon{
		HookEventName:  "afterFileEdit",
		ConversationID: "conv-1",
		GenerationID:   "gen-1",
		Model:          "claude-sonnet-4",
		CursorVersion:  "1.7.2",
		WorkspaceRoots: []string{"/repo", "/repo-2"},
		UserEmail:      ptr("dev@example.com"),
		TranscriptPath: ptr("/tmp/transcript.jsonl"),
	}
	data := cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields]{
		Common:        hookData,
		EventSpecific: &cursorcore.AfterFileEditFields{}, // event branch; common fields still tested
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}

	want := map[string]string{
		"CONVERSATION_ID": "conv-1",
		"GENERATION_ID":   "gen-1",
		"MODEL":           "claude-sonnet-4",
		"HOOK_EVENT_NAME": "afterFileEdit",
		"CURSOR_VERSION":  "1.7.2",
		"WORKSPACE_ROOTS": "/repo;/repo-2",
		"USER_EMAIL":      "dev@example.com",
		"TRANSCRIPT_PATH": "/tmp/transcript.jsonl",
	}
	for key, wantVal := range want {
		t.Run(key, func(t *testing.T) {
			assertTemplateBindingValue(t, bindings, key, wantVal)
		})
	}
}

func TestTemplateBindingsCommon_unknownKey(t *testing.T) {
	data := cursorcore.CursorHookRunData[struct{}]{
		Common: cursorcore.HookDataCommon{HookEventName: "sessionEnd"},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	_, ok := bindings.TemplateValue("SESSION_ID")
	if ok {
		t.Fatal("SESSION_ID must not be a defined placeholder")
	}
}

func TestTemplateBindingsCommon_nullOptionalJSONStillDefined(t *testing.T) {
	data := cursorcore.CursorHookRunData[struct{}]{
		Common: cursorcore.HookDataCommon{
			HookEventName:  "sessionEnd",
			ConversationID: "c",
			UserEmail:      nil,
			TranscriptPath: nil,
		},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	assertTemplateBindingValue(t, bindings, "USER_EMAIL", "")
	assertTemplateBindingValue(t, bindings, "TRANSCRIPT_PATH", "")
}

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
	data := cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields]{
		Common: cursorcore.HookDataCommon{
			HookEventName:  "afterFileEdit",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursorcore.AfterFileEditFields{FilePath: "D:/repo/file.go"},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterFileEdit")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "FILE_PATH", "D:/repo/file.go")
}

func TestTemplateBindingsAfterFileEdit_editsPlaceholderNotDefined(t *testing.T) {
	data := cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields]{
		Common: cursorcore.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursorcore.AfterFileEditFields{
			FilePath: "x",
			Edits:    []cursorcore.AfterFileEditEditPair{{OldString: "a", NewString: "b"}},
		},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	for _, key := range []string{"EDITS", "OLD_STRING", "NEW_STRING"} {
		_, ok := bindings.TemplateValue(key)
		if ok {
			t.Fatalf("%q must not be a defined placeholder", key)
		}
	}
}

func TestTemplateBindingsAfterFileEdit_unknownKey(t *testing.T) {
	data := cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields]{
		Common:        cursorcore.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursorcore.AfterFileEditFields{FilePath: "x"},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}

func TestAfterShellExecutionPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"COMMAND":  {},
		"OUTPUT":   {},
		"DURATION": {},
		"SANDBOX":  {},
	}
	if len(afterShellExecutionPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(afterShellExecutionPlaceholderExtractors))
	}
	for placeholderKey := range wantKeys {
		if _, ok := afterShellExecutionPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsAfterShellExecution_templateValueEventAndCommonFields(t *testing.T) {
	data := cursorcore.CursorHookRunData[cursorcore.AfterShellExecutionFields]{
		Common: cursorcore.HookDataCommon{
			HookEventName:  "afterShellExecution",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursorcore.AfterShellExecutionFields{
			Command:  "go test ./...",
			Output:   "PASS",
			Duration: 1234,
			Sandbox:  true,
		},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterShellExecution")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "COMMAND", "go test ./...")
	assertTemplateBindingValue(t, bindings, "OUTPUT", "PASS")
	assertTemplateBindingValue(t, bindings, "DURATION", "1234")
	assertTemplateBindingValue(t, bindings, "SANDBOX", "true")
}

func TestTemplateBindingsAfterShellExecution_decimalDuration(t *testing.T) {
	data := cursorcore.CursorHookRunData[cursorcore.AfterShellExecutionFields]{
		Common: cursorcore.HookDataCommon{HookEventName: "afterShellExecution"},
		EventSpecific: &cursorcore.AfterShellExecutionFields{
			Command:  "go test ./...",
			Output:   "PASS",
			Duration: 2841.805,
			Sandbox:  false,
		},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	assertTemplateBindingValue(t, bindings, "DURATION", "2841.805")
}

func TestTemplateBindingsAfterShellExecution_unknownKey(t *testing.T) {
	data := cursorcore.CursorHookRunData[cursorcore.AfterShellExecutionFields]{
		Common:        cursorcore.HookDataCommon{HookEventName: "afterShellExecution"},
		EventSpecific: &cursorcore.AfterShellExecutionFields{},
	}
	bindings, err := templateBindingsForCursor(&data)
	if err != nil {
		t.Fatalf("templateBindingsForCursor: %v", err)
	}
	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}

func ptr(s string) *string { return &s }

func assertTemplateBindingValue(t *testing.T, bindings template.TemplateBindings, key, want string) {
	t.Helper()
	bindingValue, ok := bindings.TemplateValue(key)
	if !ok {
		t.Fatalf("TemplateValue(%q): expected ok true", key)
	}
	if bindingValue != want {
		t.Fatalf("TemplateValue(%q): want %q, got %q", key, want, bindingValue)
	}
}
