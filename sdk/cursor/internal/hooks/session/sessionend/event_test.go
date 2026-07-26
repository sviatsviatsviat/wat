package sessionend

import (
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

func TestDecode_SessionEnd(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"sessionEnd",
		"conversation_id":"c1",
		"session_id":"s1",
		"reason":"error",
		"duration_ms":45000,
		"is_background_agent":true,
		"final_status":"failed",
		"error_message":"composer crashed"
	}`)

	if e.ConversationID != "c1" {
		t.Errorf("ConversationID = %q, want c1", e.ConversationID)
	}
	if e.SessionID != "s1" {
		t.Errorf("SessionID = %q, want s1", e.SessionID)
	}
	if e.Reason != "error" {
		t.Errorf("Reason = %q, want error", e.Reason)
	}
	if e.DurationMs != 45000 {
		t.Errorf("DurationMs = %d, want 45000", e.DurationMs)
	}
	if !e.IsBackgroundAgent {
		t.Error("IsBackgroundAgent = false, want true")
	}
	if e.FinalStatus != "failed" {
		t.Errorf("FinalStatus = %q, want failed", e.FinalStatus)
	}
	if e.ErrorMessage != "composer crashed" {
		t.Errorf("ErrorMessage = %q, want %q", e.ErrorMessage, "composer crashed")
	}
	if e.EventName() != "sessionEnd" {
		t.Errorf("EventName = %q, want sessionEnd", e.EventName())
	}
}

func TestDecode_SessionEnd_optionalErrorMessage(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"sessionEnd",
		"conversation_id":"c1",
		"reason":"completed",
		"duration_ms":1200,
		"is_background_agent":false,
		"final_status":"ok"
	}`)

	if e.Reason != "completed" {
		t.Errorf("Reason = %q, want completed", e.Reason)
	}
	if e.DurationMs != 1200 {
		t.Errorf("DurationMs = %d, want 1200", e.DurationMs)
	}
	if e.IsBackgroundAgent {
		t.Error("IsBackgroundAgent = true, want false")
	}
	if e.FinalStatus != "ok" {
		t.Errorf("FinalStatus = %q, want ok", e.FinalStatus)
	}
	if e.ErrorMessage != "" {
		t.Errorf("ErrorMessage = %q, want empty", e.ErrorMessage)
	}
}

func init() {
	register(testCodec)
}
