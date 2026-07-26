package precompact

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestDecode_PreCompact(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "pre_compact.json"))
	if err != nil {
		t.Fatal(err)
	}
	ev := mustDecode[Event](t, string(raw))

	if ev.EventName() != "preCompact" {
		t.Errorf("EventName = %q, want preCompact", ev.EventName())
	}
	if ev.ConversationID != "c1" {
		t.Errorf("ConversationID = %q, want c1", ev.ConversationID)
	}
	if ev.Trigger != "auto" {
		t.Errorf("Trigger = %q, want auto", ev.Trigger)
	}
	if ev.ContextUsagePercent != 85 {
		t.Errorf("ContextUsagePercent = %d, want 85", ev.ContextUsagePercent)
	}
	if ev.ContextTokens != 120000 {
		t.Errorf("ContextTokens = %d, want 120000", ev.ContextTokens)
	}
	if ev.ContextWindowSize != 128000 {
		t.Errorf("ContextWindowSize = %d, want 128000", ev.ContextWindowSize)
	}
	if ev.MessageCount != 45 {
		t.Errorf("MessageCount = %d, want 45", ev.MessageCount)
	}
	if ev.MessagesToCompact != 30 {
		t.Errorf("MessagesToCompact = %d, want 30", ev.MessagesToCompact)
	}
	if !ev.IsFirstCompaction {
		t.Error("IsFirstCompaction = false, want true")
	}
}

func TestDecode_PreCompact_manualTrigger(t *testing.T) {
	ev := mustDecode[Event](t, `{"hook_event_name":"preCompact","conversation_id":"c1","trigger":"manual","context_usage_percent":10,"context_tokens":1000,"context_window_size":128000,"message_count":3,"messages_to_compact":2,"is_first_compaction":false}`)
	if ev.Trigger != "manual" {
		t.Errorf("Trigger = %q, want manual", ev.Trigger)
	}
	if ev.IsFirstCompaction {
		t.Error("IsFirstCompaction = true, want false")
	}
}

func TestEncode_PreCompact_userMessage(t *testing.T) {
	out, code, err := results{}.UserMessage("Compacting context…").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"user_message":"Compacting context…"`) {
		t.Fatalf("missing user_message: %s", got)
	}
	var payload map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload) != 1 {
		t.Fatalf("payload keys = %v, want only user_message", payload)
	}
}

func TestMerge_PreCompact_userMessageLastWins(t *testing.T) {
	a := results{}.UserMessage("first")
	b := results{}.UserMessage("second")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 {
		t.Fatalf("warnings = %v, want one replacement warning", warnings)
	}
	out, code, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(string(out), `"user_message":"second"`) {
		t.Fatalf("merged user_message = %s, want second", out)
	}
}

func init() {
	register(testCodec)
}
