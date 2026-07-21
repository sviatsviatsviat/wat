package copilot

import (
	"testing"
)

func TestDecode_SessionEnd(t *testing.T) {
	e := mustDecode[SessionEnd](t, `{"hook_event_name":"SessionEnd","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","reason":"complete"}`, EventSessionEnd)
	if e.Reason != "complete" {
		t.Fatalf("SessionEnd=%+v", e)
	}
}
