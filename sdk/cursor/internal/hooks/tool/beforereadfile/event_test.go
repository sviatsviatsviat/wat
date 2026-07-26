package beforereadfile

import (
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

func TestDecode_BeforeReadFile(t *testing.T) {
	mustDecode[Event](t, `{"hook_event_name":"beforeReadFile","conversation_id":"c1","file_path":"a.go","content":"package main"}`)
}

func TestEncode_BeforeReadFileDeny_userMessageExitZero(t *testing.T) {
	out, code, err := results{}.Deny("sensitive file blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"permission":"deny"`) {
		t.Fatalf("missing permission deny: %s", got)
	}
	if !strings.Contains(got, `"user_message":"sensitive file blocked"`) {
		t.Fatalf("missing user_message: %s", got)
	}
	if strings.Contains(got, "agent_message") {
		t.Fatalf("beforeReadFile deny must not emit agent_message: %s", got)
	}
}

func TestEncode_BeforeReadFileAsk_coercesToDenyUserMessage(t *testing.T) {
	out, code, err := results{}.Ask("confirm read").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"permission":"deny"`) {
		t.Fatalf("Ask must encode as deny: %s", got)
	}
	if !strings.Contains(got, `"user_message":"confirm read"`) {
		t.Fatalf("missing user_message: %s", got)
	}
	if strings.Contains(got, `"permission":"ask"`) {
		t.Fatalf("beforeReadFile must not emit ask: %s", got)
	}
	if strings.Contains(got, "agent_message") {
		t.Fatalf("beforeReadFile Ask must not emit agent_message: %s", got)
	}
}

func init() {
	register(testCodec)
}
