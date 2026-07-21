package claude

import (
	"testing"
)

func TestDecode_TaskCreated(t *testing.T) {
	mustDecode[TaskCreated](t, `{"session_id":"s","hook_event_name":"TaskCreated","task":{"id":"t1"}}`, EventTaskCreated)
}
