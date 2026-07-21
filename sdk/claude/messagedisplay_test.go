package claude

import (
	"testing"
)

func TestDecode_MessageDisplay(t *testing.T) {
	mustDecode[MessageDisplay](t, `{"session_id":"s","hook_event_name":"MessageDisplay","turn_id":"1","message_id":"m1","index":0,"final":true,"delta":"hi"}`, EventMessageDisplay)
}
