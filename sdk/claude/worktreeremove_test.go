package claude

import (
	"testing"
)

func TestDecode_WorktreeRemove(t *testing.T) {
	mustDecode[WorktreeRemove](t, `{"session_id":"s","hook_event_name":"WorktreeRemove","worktree_path":"/wt"}`, EventWorktreeRemove)
}
