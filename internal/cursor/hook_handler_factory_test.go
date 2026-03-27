package cursor

import (
	"strings"
	"testing"
)

func newTestHookHandlerFactory() HookHandlerFactory {
	return NewHookHandlerFactory()
}

func TestHookHandlerFactory_documentedEvents(t *testing.T) {
	want := []string{
		"afterAgentResponse",
		"afterAgentThought",
		"afterFileEdit",
		"afterMCPExecution",
		"afterShellExecution",
		"afterTabFileEdit",
		"sessionEnd",
	}
	if len(cursorHookHandlerBuilders) != len(want) {
		t.Fatalf("factory has %d keys, want %d", len(cursorHookHandlerBuilders), len(want))
	}
	factory := newTestHookHandlerFactory()
	for _, eventName := range want {
		t.Run(eventName, func(t *testing.T) {
			body := []byte(`{"hook_event_name":"` + eventName + `"}`)
			handler, err := factory.HookHandlerFromJSON(body)
			if err != nil {
				t.Fatalf("HookHandlerFromJSON: %v", err)
			}
			if handler == nil {
				t.Fatal("HookHandlerFromJSON: want handler, got nil")
			}
		})
	}
}

func TestHookHandlerFactory_unknownEvent(t *testing.T) {
	factory := newTestHookHandlerFactory()
	_, err := factory.HookHandlerFromJSON([]byte(`{"hook_event_name":"preToolUse"}`))
	if err == nil {
		t.Fatal("expected error for unsupported event")
	}
	if !strings.Contains(err.Error(), "preToolUse") {
		t.Fatalf("error should name unsupported event: %v", err)
	}
}

func TestHookHandlerFactory_emptyStdinJSON(t *testing.T) {
	factory := newTestHookHandlerFactory()
	for _, body := range [][]byte{nil, {}} {
		_, err := factory.HookHandlerFromJSON(body)
		if err == nil {
			t.Fatalf("HookHandlerFromJSON(%v): want error", body)
		}
	}
}

func TestHookHandlerFactory_invalidJSON(t *testing.T) {
	factory := newTestHookHandlerFactory()
	_, err := factory.HookHandlerFromJSON([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid cursor hook JSON") {
		t.Fatalf("error should wrap unmarshal failure: %v", err)
	}
}
