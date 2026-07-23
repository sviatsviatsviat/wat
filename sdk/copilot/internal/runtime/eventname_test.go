package runtime

import (
	"errors"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"testing"
)

func TestCanonicalEventName(t *testing.T) {
	name, ok := CanonicalEventName("SessionStart")
	if !ok || name != "SessionStart" {
		t.Fatalf("CanonicalEventName(SessionStart) = %q %v", name, ok)
	}
	name, ok = CanonicalEventName("")
	if ok || name != "" {
		t.Fatalf("empty = %q %v", name, ok)
	}
	name, ok = CanonicalEventName("FutureEvent")
	if ok || name != "FutureEvent" {
		t.Fatalf("unknown = %q %v", name, ok)
	}
}

func TestEventAliasMap(t *testing.T) {
	m := EventAliasMap()
	if m["SessionStart"] != "SessionStart" || m["Stop"] != "Stop" {
		t.Fatalf("map = %v", m)
	}
	m["SessionStart"] = "mutated"
	m2 := EventAliasMap()
	if m2["SessionStart"] != "SessionStart" {
		t.Fatal("EventAliasMap should return a copy")
	}
}

func TestEventNameFromRaw(t *testing.T) {
	name, err := hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired).EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s","timestamp":"2026-01-01T00:00:00Z"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired).EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
