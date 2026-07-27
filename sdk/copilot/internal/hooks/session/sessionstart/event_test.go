package sessionstart

import (
	"context"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := results{}.Context("project uses go test ./...").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additional_context"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_SessionStart(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`))
	if err != nil {
		t.Fatal(err)
	}
	e := ev.(Event)
	if e.Source != "new" || e.InitialPrompt() != "go" {
		t.Fatalf("SessionStart=%+v", e)
	}
}

func TestMerge_SessionStart_contextJoins(t *testing.T) {
	a := results{}.Context("one")
	b := results{}.Context("two")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	raw, code, err := merged.Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(raw), `"additional_context":"one\n\ntwo"`) {
		t.Fatalf("encoded = %s", raw)
	}
	if merged.Stop() {
		t.Fatal("context should not stop")
	}
}

func TestRegisterHandler_RegistersDecoder(t *testing.T) {
	c := hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)
	d := hookkit.NewDialect(c)
	raw := []byte(`{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new"}`)

	_, err := c.Decode(raw)
	if err == nil {
		t.Fatal("expected unknown event before RegisterHandler")
	}

	RegisterHandler(d, func(context.Context, Event, Results) (Output, error) {
		return results{}.Noop(), nil
	})

	ev, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != "SessionStart" {
		t.Fatalf("EventName = %q", ev.EventName())
	}

	RegisterHandler(d, nil)
	if len(d.HandlersFor("SessionStart")) != 1 {
		t.Fatalf("nil fn should not append handlers, got %d", len(d.HandlersFor("SessionStart")))
	}
}

func init() {
	register(testCodec)
}
