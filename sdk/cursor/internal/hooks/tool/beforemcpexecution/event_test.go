package beforemcpexecution

import (
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestDecode_BeforeMCPExecution(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`)
	if e.URL != "https://mcp.example/mcp" {
		t.Fatalf("URL=%q", e.URL)
	}
}

func init() {
	Register(runtime.Codec)
}
