package cursor

import (
	"testing"
)

func TestDecode_SubagentStop(t *testing.T) {
	mustDecode[SubagentStop](t, `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","loop_count":2,"status":"completed"}`)
}
