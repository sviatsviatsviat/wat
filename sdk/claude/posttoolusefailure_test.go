package claude

import (
	"testing"
)

func TestDecode_PostToolUseFailure(t *testing.T) {
	mustDecode[PostToolUseFailure](t, `{"session_id":"s","hook_event_name":"PostToolUseFailure","tool_name":"Bash","error":"timeout"}`, EventPostToolUseFailure)
}
