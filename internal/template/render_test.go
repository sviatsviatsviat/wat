package template

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// fakeBindings marks which keys exist for the host; values may be empty when defined.
type fakeBindings struct {
	defined map[string]struct{}
	values  map[string]string
}

func (fake fakeBindings) TemplateValue(key string) (string, bool) {
	if fake.defined == nil {
		return "", false
	}
	if _, ok := fake.defined[key]; !ok {
		return "", false
	}
	if fake.values == nil {
		return "", true
	}
	return fake.values[key], true
}

func TestRenderTokens_ReplacesValues(t *testing.T) {
	argvTokens := []string{"echo", "__HOOK_EVENT_NAME__", "__CONVERSATION_ID__"}
	bindings := fakeBindings{
		defined: map[string]struct{}{
			"HOOK_EVENT_NAME": {}, "CONVERSATION_ID": {},
		},
		values: map[string]string{
			"HOOK_EVENT_NAME": "afterFileEdit",
			"CONVERSATION_ID": "conv-9",
		},
	}

	renderedArgv, unknownPlaceholders := RenderTokens(argvTokens, bindings)
	if len(unknownPlaceholders) != 0 {
		t.Fatalf("unknownPlaceholders = %v", unknownPlaceholders)
	}
	if renderedArgv[0] != "echo" {
		t.Fatalf("literal token passthrough: want echo, got %q", renderedArgv[0])
	}
	if renderedArgv[1] != "afterFileEdit" || renderedArgv[2] != "conv-9" {
		t.Fatalf("unexpected rendered tokens: %v", renderedArgv)
	}
}

func TestRenderTokens_UnknownPlaceholder(t *testing.T) {
	argvTokens := []string{"echo", "__DOES_NOT_EXIST__", "__USER_EMAIL__"}
	bindings := fakeBindings{
		defined: map[string]struct{}{"USER_EMAIL": {}},
		values:  map[string]string{"USER_EMAIL": ""},
	}

	renderedArgv, unknownPlaceholders := RenderTokens(argvTokens, bindings)
	if len(unknownPlaceholders) != 1 || unknownPlaceholders[0] != "DOES_NOT_EXIST" {
		t.Fatalf("unexpected unknownPlaceholders: %v", unknownPlaceholders)
	}
	// Undefined placeholders substitute to empty string in the argv token (see RenderTokens).
	if renderedArgv[1] != "" {
		t.Fatalf("unknown placeholder token should render empty, got %q", renderedArgv[1])
	}
	if renderedArgv[2] != "" {
		t.Fatalf("empty defined placeholder should substitute empty string, got %q", renderedArgv[2])
	}
}

func TestRenderTokens_UndefinedKeyEvenWithValueMap(t *testing.T) {
	argvTokens := []string{"__SECRET__"}
	bindings := fakeBindings{
		defined: map[string]struct{}{"OTHER": {}},
		values:  map[string]string{"SECRET": "x"},
	}
	_, unknownPlaceholders := RenderTokens(argvTokens, bindings)
	if len(unknownPlaceholders) != 1 || unknownPlaceholders[0] != "SECRET" {
		t.Fatalf("want unknown SECRET, got %v", unknownPlaceholders)
	}
}

var _ core.TemplateBindings = fakeBindings{}
