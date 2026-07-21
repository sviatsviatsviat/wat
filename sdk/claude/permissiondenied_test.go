package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_PermissionDenied(t *testing.T) {
	mustDecode[PermissionDenied](t, `{"session_id":"s","hook_event_name":"PermissionDenied","tool_name":"Bash"}`, EventPermissionDenied)

	if !(permissionDeniedOutput{}).IsZero() {
		t.Fatal("empty permission denied output should be zero")
	}
	retry := permissionDeniedResults{}.Retry()
	if retry.IsZero() {
		t.Fatal("Retry() should be non-zero")
	}

	out, code, err := retry.Encode()
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
	if hso["hookEventName"] != EventPermissionDenied {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventPermissionDenied)
	}
	if hso["retry"] != true {
		t.Fatalf("retry = %v", hso["retry"])
	}
}
