package claude

import (
	"testing"
)

func TestDecode_Notification(t *testing.T) {
	mustDecode[Notification](t, `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`, EventNotification)
}
