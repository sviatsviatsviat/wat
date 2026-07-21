package claude

import (
	"testing"
)

func TestDecode_TeammateIdle(t *testing.T) {
	mustDecode[TeammateIdle](t, `{"session_id":"s","hook_event_name":"TeammateIdle"}`, EventTeammateIdle)
}
