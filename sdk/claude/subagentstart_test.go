package claude

import (
	"testing"
)

func TestDecode_SubagentStart(t *testing.T) {
	mustDecode[SubagentStart](t, `{"session_id":"s","hook_event_name":"SubagentStart","agent_id":"a1","agent_type":"reviewer"}`, EventSubagentStart)
}
