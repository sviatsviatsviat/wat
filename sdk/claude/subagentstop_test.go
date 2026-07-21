package claude

import (
	"testing"
)

func TestDecode_SubagentStop(t *testing.T) {
	mustDecode[SubagentStop](t, `{"session_id":"s","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"worker","stop_hook_active":true,"last_assistant_message":"done"}`, EventSubagentStop)
}
