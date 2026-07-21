package claude

import (
	"testing"
)

func TestDecode_CwdChanged(t *testing.T) {
	mustDecode[CwdChanged](t, `{"session_id":"s","hook_event_name":"CwdChanged","new_cwd":"/new","old_cwd":"/old"}`, EventCwdChanged)
}
