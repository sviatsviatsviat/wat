package claude

import (
	"errors"
	"testing"
)

func TestEventNameFromRaw(t *testing.T) {
	name, err := eventNameFromRaw([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = eventNameFromRaw([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEmptyPayload) && err.Error() == "" {
		t.Fatalf("err = %v", err)
	}
}
