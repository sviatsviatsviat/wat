package cursor

import (
	"testing"
)

func TestDecode_SubagentStart(t *testing.T) {
	mustDecode[SubagentStart](t, `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files"}`)
}
