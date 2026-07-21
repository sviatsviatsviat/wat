package copilot

import (
	"testing"
)

func TestDecode_SubagentStop(t *testing.T) {
	e := mustDecode[SubagentStop](t, `{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","agent_display_name":"Task","stop_reason":"end_turn"}`, EventSubagentStop)
	if e.Name() != "task" || e.Reason() != "end_turn" {
		t.Fatalf("SubagentStop=%+v", e)
	}
}
