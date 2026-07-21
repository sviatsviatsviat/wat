package copilot

import (
	"testing"
)

func TestDecode_SubagentStart(t *testing.T) {
	e := mustDecode[SubagentStart](t, `{"hook_event_name":"SubagentStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"explore","agent_display_name":"Explore","agent_description":"search codebase"}`, EventSubagentStart)
	if e.Name() != "explore" || e.DisplayName() != "Explore" {
		t.Fatalf("SubagentStart=%+v", e)
	}
}
