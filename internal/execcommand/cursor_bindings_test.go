package execcommand

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func TestTemplateBindingsCommon_templateValueAllCommonFields(t *testing.T) {
	hookData := cursor.HookDataCommon{
		HookEventName:  "afterFileEdit",
		ConversationID: "conv-1",
		GenerationID:   "gen-1",
		Model:          "claude-sonnet-4",
		CursorVersion:  "1.7.2",
		WorkspaceRoots: []string{"/repo", "/repo-2"},
		UserEmail:      ptr("dev@example.com"),
		TranscriptPath: ptr("/tmp/transcript.jsonl"),
	}
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common:        hookData,
		EventSpecific: &cursor.AfterFileEditFields{}, // event branch; common fields still tested
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterFileEditPlaceholderExtractors)

	want := map[string]string{
		"CONVERSATION_ID": "conv-1",
		"GENERATION_ID":   "gen-1",
		"MODEL":           "claude-sonnet-4",
		"HOOK_EVENT_NAME": "afterFileEdit",
		"CURSOR_VERSION":  "1.7.2",
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
	data := cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	}
	bindings := newTemplateBindingsCommon(data.Common)
	_, ok := bindings.TemplateValue("SESSION_ID")
	if ok {
		t.Fatal("SESSION_ID must not be a defined placeholder")
	}
}

func TestTemplateBindingsCommon_nullOptionalJSONStillDefined(t *testing.T) {
	data := cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{
			HookEventName:  "sessionEnd",
			ConversationID: "c",
			UserEmail:      nil,
			TranscriptPath: nil,
		},
	}
	bindings := newTemplateBindingsCommon(data.Common)
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
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "afterFileEdit",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursor.AfterFileEditFields{FilePath: "D:/repo/file.go"},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterFileEditPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterFileEdit")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "FILE_PATH", "D:/repo/file.go")
}

func TestTemplateBindingsAfterFileEdit_editsPlaceholderNotDefined(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{
			FilePath: "x",
			Edits:    []cursor.AfterFileEditEditPair{{OldString: "a", NewString: "b"}},
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterFileEditPlaceholderExtractors)
	for _, key := range []string{"EDITS", "OLD_STRING", "NEW_STRING"} {
		_, ok := bindings.TemplateValue(key)
		if ok {
			t.Fatalf("%q must not be a defined placeholder", key)
		}
	}
}

func TestTemplateBindingsAfterFileEdit_unknownKey(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common:        cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{FilePath: "x"},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterFileEditPlaceholderExtractors)
	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}

func TestAfterShellExecutionPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
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
	data := cursor.CursorHookRunData[cursor.AfterShellExecutionFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "afterShellExecution",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursor.AfterShellExecutionFields{
			Command:  "go test ./...",
			Output:   "PASS",
			Duration: 1234,
			Sandbox:  true,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterShellExecutionPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterShellExecution")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "DURATION", "1234")
	assertTemplateBindingValue(t, bindings, "SANDBOX", "true")
}

func TestTemplateBindingsAfterShellExecution_decimalDuration(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterShellExecutionFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterShellExecution"},
		EventSpecific: &cursor.AfterShellExecutionFields{
			Command:  "go test ./...",
			Output:   "PASS",
			Duration: 2841.805,
			Sandbox:  false,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterShellExecutionPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "DURATION", "2841.805")
}

func TestTemplateBindingsAfterShellExecution_unknownKey(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterShellExecutionFields]{
		Common:        cursor.HookDataCommon{HookEventName: "afterShellExecution"},
		EventSpecific: &cursor.AfterShellExecutionFields{},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterShellExecutionPlaceholderExtractors)
	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}

func TestAfterMCPExecutionPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"TOOL_NAME": {},
		"DURATION":  {},
	}
	if len(afterMCPExecutionPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(afterMCPExecutionPlaceholderExtractors))
	}
	for placeholderKey := range wantKeys {
		if _, ok := afterMCPExecutionPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsAfterMCPExecution_templateValueEventAndCommonFields(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterMCPExecutionFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "afterMCPExecution",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursor.AfterMCPExecutionFields{
			ToolName:   "search",
			ToolInput:  `{"q":"x"}`,
			ResultJSON: `{"hits":[]}`,
			Duration:   1234,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterMCPExecutionPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterMCPExecution")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "TOOL_NAME", "search")
	assertTemplateBindingValue(t, bindings, "DURATION", "1234")
}

func TestTemplateBindingsAfterMCPExecution_decimalDuration(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterMCPExecutionFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterMCPExecution"},
		EventSpecific: &cursor.AfterMCPExecutionFields{
			ToolName: "t",
			Duration: 2841.805,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterMCPExecutionPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "DURATION", "2841.805")
}

func TestAfterAgentThoughtPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"DURATION_MS": {},
	}
	if len(afterAgentThoughtPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(afterAgentThoughtPlaceholderExtractors))
	}
	for placeholderKey := range wantKeys {
		if _, ok := afterAgentThoughtPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsAfterAgentThought_templateValueEventAndCommonFields(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterAgentThoughtFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "afterAgentThought",
			ConversationID: "conv-1",
		},
		EventSpecific: &cursor.AfterAgentThoughtFields{
			Text:       "step by step",
			DurationMs: 5000,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterAgentThoughtPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterAgentThought")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "DURATION_MS", "5000")
}

func TestTemplateBindingsAfterAgentThought_textNotAPlaceholder(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterAgentThoughtFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterAgentThought"},
		EventSpecific: &cursor.AfterAgentThoughtFields{
			Text:       "secret",
			DurationMs: 1,
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterAgentThoughtPlaceholderExtractors)
	_, ok := bindings.TemplateValue("TEXT")
	if ok {
		t.Fatal("TEXT must not be a defined placeholder")
	}
}

func TestSessionEndPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"SESSION_ID": {}, "REASON": {}, "DURATION_MS": {}, "IS_BACKGROUND": {},
		"FINAL_STATUS": {}, "ERROR_MESSAGE": {},
	}
	if len(sessionEndPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(sessionEndPlaceholderExtractors))
	}
	for placeholderKey := range wantKeys {
		if _, ok := sessionEndPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsSessionEnd_templateValueEventAndCommonFields(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.SessionEndFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "sessionEnd",
			ConversationID: "conv-se",
		},
		EventSpecific: &cursor.SessionEndFields{
			SessionID:         "sid-42",
			Reason:            "error",
			DurationMs:        45000,
			IsBackgroundAgent: true,
			FinalStatus:       "done",
			ErrorMessage:      "oops",
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, sessionEndPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "sessionEnd")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-se")
	assertTemplateBindingValue(t, bindings, "SESSION_ID", "sid-42")
	assertTemplateBindingValue(t, bindings, "REASON", "error")
	assertTemplateBindingValue(t, bindings, "DURATION_MS", "45000")
	assertTemplateBindingValue(t, bindings, "IS_BACKGROUND", "true")
	assertTemplateBindingValue(t, bindings, "FINAL_STATUS", "done")
	assertTemplateBindingValue(t, bindings, "ERROR_MESSAGE", "oops")
}

func TestTemplateBindingsSessionEnd_emptyOptionalErrorMessage(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.SessionEndFields]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
		EventSpecific: &cursor.SessionEndFields{
			SessionID:   "s",
			Reason:      "completed",
			DurationMs:  1,
			FinalStatus: "ok",
		},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, sessionEndPlaceholderExtractors)
	assertTemplateBindingValue(t, bindings, "ERROR_MESSAGE", "")
}

func TestTemplateBindingsAfterMCPExecution_unknownKey(t *testing.T) {
	data := cursor.CursorHookRunData[cursor.AfterMCPExecutionFields]{
		Common:        cursor.HookDataCommon{HookEventName: "afterMCPExecution"},
		EventSpecific: &cursor.AfterMCPExecutionFields{ToolName: "x"},
	}
	bindings := templateBindingsFromCursorEventPayload(data.Common, data.EventSpecific, afterMCPExecutionPlaceholderExtractors)
	for _, key := range []string{"SANDBOX", "TOOL_INPUT", "RESULT_JSON"} {
		_, ok := bindings.TemplateValue(key)
		if ok {
			t.Fatalf("%q must not be a defined placeholder", key)
		}
	}
}

func ptr(s string) *string { return &s }

func assertTemplateBindingValue(t *testing.T, bindings templateBindings, key, want string) {
	t.Helper()
	bindingValue, ok := bindings.TemplateValue(key)
	if !ok {
		t.Fatalf("TemplateValue(%q): expected ok true", key)
	}
	if bindingValue != want {
		t.Fatalf("TemplateValue(%q): want %q, got %q", key, want, bindingValue)
	}
}
