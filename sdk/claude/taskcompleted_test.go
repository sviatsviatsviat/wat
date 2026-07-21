package claude

import (
	"testing"
)

func TestDecode_TaskCompleted(t *testing.T) {
	mustDecode[TaskCompleted](t, `{"session_id":"s","hook_event_name":"TaskCompleted","task":{"id":"t1"}}`, EventTaskCompleted)
}
