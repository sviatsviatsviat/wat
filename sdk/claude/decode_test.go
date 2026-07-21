package claude

import (
	"errors"
	"strings"
	"testing"
)

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := codec.Decode([]byte(raw))
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

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := codec.Decode([]byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := codec.Decode([]byte(`{"session_id":"s1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	t.Run("envelope", func(t *testing.T) {
		_, err := codec.Decode([]byte("not json"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("typed event", func(t *testing.T) {
		raw := []byte(`{"session_id":"s1","hook_event_name":"PreToolUse","cwd":"/w","tool_name":123}`)
		_, err := codec.Decode(raw)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrDecodePayload) {
			t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
		}
		if !strings.Contains(err.Error(), "decode payload") {
			t.Fatalf("error = %v", err)
		}
	})
}
