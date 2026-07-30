package elicitation

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_Elicitation(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"Elicitation",
		"mcp_server_name":"my-mcp-server",
		"message":"Please provide your credentials",
		"mode":"form",
		"elicitation_id":"elicit-123",
		"requested_schema":{"type":"object","properties":{"username":{"type":"string"}}}
	}`, event.Elicitation)
	if ev.MCPServerName != "my-mcp-server" {
		t.Fatalf("MCPServerName = %q", ev.MCPServerName)
	}
	if ev.Message != "Please provide your credentials" {
		t.Fatalf("Message = %q", ev.Message)
	}
	if ev.Mode != "form" {
		t.Fatalf("Mode = %q", ev.Mode)
	}
	if ev.ElicitationID != "elicit-123" {
		t.Fatalf("ElicitationID = %q", ev.ElicitationID)
	}
	if len(ev.RequestedSchema) == 0 {
		t.Fatal("RequestedSchema is empty")
	}

	urlEv := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"Elicitation",
		"mcp_server_name":"my-mcp-server",
		"message":"Please authenticate",
		"mode":"url",
		"url":"https://auth.example.com/login"
	}`, event.Elicitation)
	if urlEv.Mode != "url" {
		t.Fatalf("Mode = %q", urlEv.Mode)
	}
	if urlEv.URL != "https://auth.example.com/login" {
		t.Fatalf("URL = %q", urlEv.URL)
	}

	out, code, err := results{}.Accept().WithContent(map[string]any{"answer": "yes"}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso, ok := got["hookSpecificOutput"].(map[string]any)
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %s", out)
	}
	if hso["hookEventName"] != event.Elicitation {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], event.Elicitation)
	}
	if hso["action"] != "accept" {
		t.Fatalf("action = %v", hso["action"])
	}
	content, ok := hso["content"].(map[string]any)
	if !ok || content["answer"] != "yes" {
		t.Fatalf("content = %v", hso["content"])
	}

	// Results has no Noop(); empty output is the zero-state contract.
	if !(output{}).IsZero() {
		t.Fatal("empty elicitation output should be zero")
	}
}

func init() {
	register(testCodec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != wantName {
		t.Fatalf("EventName() = %q, want %q", ev.EventName(), wantName)
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}
