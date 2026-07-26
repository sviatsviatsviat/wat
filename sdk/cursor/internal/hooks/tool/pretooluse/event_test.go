package pretooluse

import (
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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

func TestToolInput_AsShell(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1"}`))
	if err != nil {
		t.Fatal(err)
	}
	pre := ev.(Event)
	input, ok := pre.ToolInput.AsShell()
	if !ok || input.Command != "ls" {
		t.Fatalf("AsShell = %+v, %v", input, ok)
	}
}

func TestDecode_PreToolUse(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":"Shell","tool_input":{"command":"ls"},"tool_use_id":"t1","agent_message":"Installing dependencies..."}`)
	if e.ShellCommand() != "ls" {
		t.Fatalf("ShellCommand=%q", e.ShellCommand())
	}
	if e.AgentMessage != "Installing dependencies..." {
		t.Fatalf("AgentMessage=%q", e.AgentMessage)
	}
}

func TestEncode_PreToolUseAsk_schemaAcceptedNotEnforced(t *testing.T) {
	out, code, err := results{}.Ask("confirm with user").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"permission":"ask"`) {
		t.Fatalf("Ask must still encode schema-accepted ask: %s", got)
	}
	if !strings.Contains(got, `"agent_message":"confirm with user"`) {
		t.Fatalf("missing agent_message: %s", got)
	}
}

func TestDecode_PreToolUse_InvalidJSON(t *testing.T) {
	raw := []byte(`{"hook_event_name":"preToolUse","conversation_id":"c1","tool_name":123}`)
	_, err := testCodec.Decode(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, runtime.ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func init() {
	register(testCodec)
}
