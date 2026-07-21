package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_WorktreeCreate(t *testing.T) {
	mustDecode[WorktreeCreate](t, `{"session_id":"s","hook_event_name":"WorktreeCreate"}`, EventWorktreeCreate)

	out, code, err := worktreeCreateResults{}.Path("/tmp/wt").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d, want %d", code, SuccessExit)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	hso, ok := got["hookSpecificOutput"].(map[string]any)
	if !ok {
		t.Fatalf("missing hookSpecificOutput: %s", out)
	}
	if hso["hookEventName"] != EventWorktreeCreate {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventWorktreeCreate)
	}
	if hso["worktreePath"] != "/tmp/wt" {
		t.Fatalf("worktreePath = %v", hso["worktreePath"])
	}

	// WorktreeCreateResults has no Noop(); empty output is the zero-state contract.
	if !(worktreeCreateOutput{}).IsZero() {
		t.Fatal("empty worktree create output should be zero")
	}
}
