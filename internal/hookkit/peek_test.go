package hookkit

import (
	"errors"
	"testing"
)

func TestPeekHookEventName(t *testing.T) {
	t.Parallel()
	got, err := PeekHookEventName([]byte(`{"hook_event_name":"PreToolUse","cwd":"/w"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != "PreToolUse" {
		t.Fatalf("PeekHookEventName = %q, want PreToolUse", got)
	}
	got, err = PeekHookEventName([]byte(`{"cwd":"/w"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("missing name = %q, want empty", got)
	}
	_, err = PeekHookEventName([]byte(`{`))
	if err == nil {
		t.Fatal("expected JSON error")
	}
}

func TestRequireHookEventName(t *testing.T) {
	t.Parallel()
	empty := errors.New("empty")
	decodeErr := errors.New("decode")
	nameRequired := errors.New("name required")

	got, err := RequireHookEventName([]byte(`{"hook_event_name":"Stop"}`), empty, decodeErr, nameRequired)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Stop" {
		t.Fatalf("got = %q", got)
	}

	_, err = RequireHookEventName(nil, empty, decodeErr, nameRequired)
	if !errors.Is(err, empty) {
		t.Fatalf("empty: %v", err)
	}

	_, err = RequireHookEventName([]byte(`{"cwd":"/w"}`), empty, decodeErr, nameRequired)
	if !errors.Is(err, nameRequired) {
		t.Fatalf("name required: %v", err)
	}

	_, err = RequireHookEventName([]byte(`{`), empty, decodeErr, nameRequired)
	if !errors.Is(err, decodeErr) {
		t.Fatalf("decode: %v", err)
	}
}
