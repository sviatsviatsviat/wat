package claude

import (
	"testing"
)

func TestDecode_PostToolUse(t *testing.T) {
	mustDecode[PostToolUse](t, `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`, EventPostToolUse)
}
