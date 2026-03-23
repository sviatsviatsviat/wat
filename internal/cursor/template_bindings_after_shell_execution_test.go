package cursor

import "testing"

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
	input := `{
		"hook_event_name": "afterShellExecution",
		"conversation_id": "conv-1",
		"command": "go test ./...",
		"output": "PASS",
		"duration": 1234,
		"sandbox": true
	}`
	commonData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution([]byte(input), commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}
	bindings := newTemplateBindingsAfterShellExecution(hookData)

	assertTemplateBindingValue(t, bindings, "HOOK_EVENT_NAME", "afterShellExecution")
	assertTemplateBindingValue(t, bindings, "CONVERSATION_ID", "conv-1")
	assertTemplateBindingValue(t, bindings, "COMMAND", "go test ./...")
	assertTemplateBindingValue(t, bindings, "OUTPUT", "PASS")
	assertTemplateBindingValue(t, bindings, "DURATION", "1234")
	assertTemplateBindingValue(t, bindings, "SANDBOX", "true")
}

func TestTemplateBindingsAfterShellExecution_decimalDuration(t *testing.T) {
	input := `{
		"hook_event_name": "afterShellExecution",
		"command": "go test ./...",
		"output": "PASS",
		"duration": 2841.805,
		"sandbox": false
	}`
	commonData, err := newHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution([]byte(input), commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}
	bindings := newTemplateBindingsAfterShellExecution(hookData)
	assertTemplateBindingValue(t, bindings, "DURATION", "2841.805")
}

func TestTemplateBindingsAfterShellExecution_unknownKey(t *testing.T) {
	raw := []byte(`{"hook_event_name":"afterShellExecution"}`)
	commonData, err := newHookDataCommon(raw)
	if err != nil {
		t.Fatalf("newHookDataCommon: %v", err)
	}
	hookData, err := newHookDataAfterShellExecution(raw, commonData)
	if err != nil {
		t.Fatalf("newHookDataAfterShellExecution: %v", err)
	}
	bindings := newTemplateBindingsAfterShellExecution(hookData)
	_, ok := bindings.TemplateValue("TOOL_NAME")
	if ok {
		t.Fatal("TOOL_NAME must not be a defined placeholder")
	}
}
