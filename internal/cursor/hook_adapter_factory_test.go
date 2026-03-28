package cursor

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func newTestHookAdapterFactory() HookAdapterFactory {
	return NewHookAdapterFactory()
}

func TestHookAdapterFactory_documentedEvents(t *testing.T) {
	want := []string{
		"afterAgentResponse",
		"afterAgentThought",
		"afterFileEdit",
		"afterMCPExecution",
		"afterShellExecution",
		"afterTabFileEdit",
		"sessionEnd",
	}
	if len(cursorHookAdapterBuilders) != len(want) {
		t.Fatalf("factory has %d keys, want %d", len(cursorHookAdapterBuilders), len(want))
	}
	factory := newTestHookAdapterFactory()
	mock := cli.NewMockConsole()
	for _, eventName := range want {
		t.Run(eventName, func(t *testing.T) {
			body := []byte(`{"hook_event_name":"` + eventName + `"}`)
			adapter, err := factory.HookAdapterFromJSON(body, mock)
			if err != nil {
				t.Fatalf("HookAdapterFromJSON: %v", err)
			}
			if adapter == nil {
				t.Fatal("HookAdapterFromJSON: want adapter, got nil")
			}
		})
	}
}

func TestHookAdapterFactory_unknownEvent(t *testing.T) {
	factory := newTestHookAdapterFactory()
	mock := cli.NewMockConsole()
	_, err := factory.HookAdapterFromJSON([]byte(`{"hook_event_name":"preToolUse"}`), mock)
	if err == nil {
		t.Fatal("expected error for unsupported event")
	}
	if !strings.Contains(err.Error(), "preToolUse") {
		t.Fatalf("error should name unsupported event: %v", err)
	}
}

func TestHookAdapterFactory_emptyStdinJSON(t *testing.T) {
	factory := newTestHookAdapterFactory()
	mock := cli.NewMockConsole()
	for _, body := range [][]byte{nil, {}} {
		_, err := factory.HookAdapterFromJSON(body, mock)
		if err == nil {
			t.Fatalf("HookAdapterFromJSON(%v): want error", body)
		}
	}
}

func TestHookAdapterFactory_invalidJSON(t *testing.T) {
	factory := newTestHookAdapterFactory()
	mock := cli.NewMockConsole()
	_, err := factory.HookAdapterFromJSON([]byte(`not json`), mock)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid cursor hook JSON") {
		t.Fatalf("error should wrap unmarshal failure: %v", err)
	}
}
