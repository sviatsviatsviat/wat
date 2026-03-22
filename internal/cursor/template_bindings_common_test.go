package cursor

import "testing"

func TestCommonPlaceholderExtractors_registry(t *testing.T) {
	wantKeys := map[string]struct{}{
		"CONVERSATION_ID": {},
		"GENERATION_ID":   {},
		"MODEL":           {},
		"HOOK_EVENT_NAME": {},
		"CURSOR_VERSION":  {},
		"WORKSPACE_ROOTS": {},
		"USER_EMAIL":      {},
		"TRANSCRIPT_PATH": {},
	}
	if len(commonPlaceholderExtractors) != len(wantKeys) {
		t.Fatalf("registry size: want %d, got %d", len(wantKeys), len(commonPlaceholderExtractors))
	}
	for placeholderKey := range commonPlaceholderExtractors {
		if _, ok := wantKeys[placeholderKey]; !ok {
			t.Fatalf("unexpected registry key %q", placeholderKey)
		}
	}
	for placeholderKey := range wantKeys {
		if _, ok := commonPlaceholderExtractors[placeholderKey]; !ok {
			t.Fatalf("missing registry key %q", placeholderKey)
		}
	}
}

func TestTemplateBindingsCommon_templateValueAllCommonFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "conv-1",
		"generation_id": "gen-1",
		"model": "claude-sonnet-4",
		"cursor_version": "1.7.2",
		"workspace_roots": ["/repo", "/repo-2"],
		"user_email": "dev@example.com",
		"transcript_path": "/tmp/transcript.jsonl"
	}`
	hookData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	bindings := newTemplateBindingsCommon(hookData)

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
	hookData, err := newHookDataCommon([]byte(`{"hook_event_name":"afterFileEdit"}`))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	bindings := newTemplateBindingsCommon(hookData)
	_, ok := bindings.TemplateValue("SESSION_ID")
	if ok {
		t.Fatal("SESSION_ID must not be a defined placeholder")
	}
}

func TestTemplateBindingsCommon_nullOptionalJSONStillDefined(t *testing.T) {
	input := `{
		"hook_event_name": "sessionEnd",
		"conversation_id": "c",
		"user_email": null,
		"transcript_path": null
	}`
	hookData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	bindings := newTemplateBindingsCommon(hookData)
	assertTemplateBindingValue(t, bindings, "USER_EMAIL", "")
	assertTemplateBindingValue(t, bindings, "TRANSCRIPT_PATH", "")
}

func assertTemplateBindingValue(t *testing.T, bindings interface {
	TemplateValue(key string) (string, bool)
}, key, want string,
) {
	t.Helper()
	bindingValue, ok := bindings.TemplateValue(key)
	if !ok {
		t.Fatalf("TemplateValue(%q): expected ok true", key)
	}
	if bindingValue != want {
		t.Fatalf("TemplateValue(%q): want %q, got %q", key, want, bindingValue)
	}
}
