package elicitationresult

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_ElicitationResult(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"ElicitationResult",
		"mcp_server_name":"my-mcp-server",
		"action":"accept",
		"content":{"username":"alice"},
		"mode":"form",
		"elicitation_id":"elicit-123"
	}`, event.ElicitationResult)
	if ev.MCPServerName != "my-mcp-server" {
		t.Fatalf("MCPServerName = %q", ev.MCPServerName)
	}
	if ev.Action != "accept" {
		t.Fatalf("Action = %q", ev.Action)
	}
	if ev.Mode != "form" {
		t.Fatalf("Mode = %q", ev.Mode)
	}
	if ev.ElicitationID != "elicit-123" {
		t.Fatalf("ElicitationID = %q", ev.ElicitationID)
	}
	if string(ev.Content) != `{"username":"alice"}` {
		t.Fatalf("Content = %s", ev.Content)
	}
}

func TestEncode_Decline(t *testing.T) {
	out, code, err := results{}.Decline().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso, ok := got["hookSpecificOutput"].(map[string]any)
	if !ok || hso["hookEventName"] != event.ElicitationResult || hso["action"] != "decline" {
		t.Fatalf("got %s", out)
	}
	if !(results{}.Decline().Stop()) {
		t.Fatal("decline should stop")
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
