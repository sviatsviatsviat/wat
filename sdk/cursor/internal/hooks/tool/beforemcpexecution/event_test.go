package beforemcpexecution

import (
	"encoding/json"
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

func TestDecode_BeforeMCPExecution(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"beforeMCPExecution","conversation_id":"c1","tool_name":"MCP:browser_navigate","tool_input":"{}","url":"https://mcp.example/mcp"}`)
	if e.URL != "https://mcp.example/mcp" {
		t.Fatalf("URL=%q", e.URL)
	}
}

func TestEncode_BeforeMCPAsk_enforcedPermission(t *testing.T) {
	out, code, err := event.NewPermissionResults().Ask("confirm MCP tool").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (ask is not deny/exit 2)", code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["permission"] != "ask" {
		t.Fatalf("permission = %v, want ask (enforced on beforeMCPExecution)", got["permission"])
	}
	if got["agent_message"] != "confirm MCP tool" {
		t.Fatalf("agent_message = %v, want confirm MCP tool", got["agent_message"])
	}
	if _, ok := got["user_message"]; ok {
		t.Fatalf("Ask default must not set user_message: %s", out)
	}
}

func init() {
	register(testCodec)
}
