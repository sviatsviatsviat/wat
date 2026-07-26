package aftermcpexecution

import (
	"context"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
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

func TestDecode_AfterMCPExecution(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":{"url":"https://example.com"},"result_json":"{\"ok\":true}","duration":42}`)
	if e.ToolName != "MCP:browser_navigate" {
		t.Fatalf("ToolName = %q", e.ToolName)
	}
	if e.ResultJSON != `{"ok":true}` {
		t.Fatalf("ResultJSON = %q", e.ResultJSON)
	}
	if e.DurationMillis() != 42 {
		t.Fatalf("DurationMillis = %d", e.DurationMillis())
	}
	if string(e.ToolInput.Raw()) == "" {
		t.Fatal("ToolInput should bind tool_input")
	}
}

func TestDurationMillis_PrefersDurationOverDurationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"afterMCPExecution",
		"conversation_id":"c1",
		"tool_name":"MCP:x",
		"result_json":"{}",
		"duration":100,
		"duration_ms":999
	}`)
	if e.DurationMillis() != 100 {
		t.Fatalf("DurationMillis()=%d, want documented duration field 100", e.DurationMillis())
	}
}

func TestDurationMillis_ExplicitZeroBeatsDurationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"afterMCPExecution",
		"conversation_id":"c1",
		"tool_name":"MCP:x",
		"result_json":"{}",
		"duration":0,
		"duration_ms":999
	}`)
	if e.DurationMillis() != 0 {
		t.Fatalf("DurationMillis()=%d, want explicit duration 0", e.DurationMillis())
	}
}

func TestDurationMillis_FallsBackToDurationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"afterMCPExecution",
		"conversation_id":"c1",
		"tool_name":"MCP:x",
		"result_json":"{}",
		"duration_ms":50
	}`)
	if e.DurationMillis() != 50 {
		t.Fatalf("DurationMillis()=%d, want 50", e.DurationMillis())
	}
}

func TestRegisterHandler_observeOnly(t *testing.T) {
	d := hookkit.NewDialect(testCodec)
	called := false
	RegisterHandler(d, func(_ context.Context, e Event) error {
		called = true
		if e.ToolName != "MCP:browser_navigate" {
			t.Fatalf("ToolName = %q", e.ToolName)
		}
		return nil
	})
	handlers := d.HandlersFor(event.AfterMCPExecution)
	if len(handlers) != 1 {
		t.Fatalf("handlers = %d, want 1", len(handlers))
	}
	ev := mustDecode[Event](t, `{"hook_event_name":"afterMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","result_json":"{}","duration":1}`)
	out, err := handlers[0].Invoke(context.Background(), ev)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("observe handler must emit no output, got %#v", out)
	}
	if !called {
		t.Fatal("handler was not invoked")
	}
}

func TestRegisterHandler_nilFn(t *testing.T) {
	d := hookkit.NewDialect(testCodec)
	RegisterHandler(d, nil)
	if got := len(d.HandlersFor(event.AfterMCPExecution)); got != 0 {
		t.Fatalf("handlers = %d, want 0", got)
	}
}

func init() {
	register(testCodec)
}
