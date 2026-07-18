package claude

import "testing"

func TestEventNameFromRaw(t *testing.T) {
	name, err := eventNameFromRaw([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`), "")
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	name, err = eventNameFromRaw([]byte(`{"session_id":"s"}`), "SessionStart")
	if err != nil {
		t.Fatal(err)
	}
	if name != "SessionStart" {
		t.Fatalf("hint name = %q", name)
	}
}
