package claude

import (
	"testing"
)

func TestDecode_SessionEnd(t *testing.T) {
	mustDecode[SessionEnd](t, `{"session_id":"s","hook_event_name":"SessionEnd","reason":"clear"}`, EventSessionEnd)
}
