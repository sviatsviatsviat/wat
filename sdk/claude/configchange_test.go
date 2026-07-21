package claude

import (
	"testing"
)

func TestDecode_ConfigChange(t *testing.T) {
	mustDecode[ConfigChange](t, `{"session_id":"s","hook_event_name":"ConfigChange","source":"user_settings"}`, EventConfigChange)
}
