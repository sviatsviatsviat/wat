package claude

import (
	"testing"
)

func TestDecode_WorktreeCreate(t *testing.T) {
	mustDecode[WorktreeCreate](t, `{"session_id":"s","hook_event_name":"WorktreeCreate"}`, EventWorktreeCreate)
}
