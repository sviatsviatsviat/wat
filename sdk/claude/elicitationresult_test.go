package claude

import (
	"testing"
)

func TestDecode_ElicitationResult(t *testing.T) {
	mustDecode[ElicitationResult](t, `{"session_id":"s","hook_event_name":"ElicitationResult","server_name":"srv","action":"accept"}`, EventElicitationResult)
}
