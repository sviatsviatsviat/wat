package claude

import (
	"testing"
)

func TestDecode_PermissionDenied(t *testing.T) {
	mustDecode[PermissionDenied](t, `{"session_id":"s","hook_event_name":"PermissionDenied","tool_name":"Bash"}`, EventPermissionDenied)
}
