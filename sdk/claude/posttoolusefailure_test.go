package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_PostToolUseFailure(t *testing.T) {
	mustDecode[PostToolUseFailure](t, `{"session_id":"s","hook_event_name":"PostToolUseFailure","tool_name":"Bash","error":"timeout"}`, EventPostToolUseFailure)
}

func TestEncode_PostToolUseFailureContext(t *testing.T) {
	out, code, err := postToolUseFailureResults{}.Context("retry with smaller input").Encode()
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
	if hso["hookEventName"] != EventPostToolUseFailure {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventPostToolUseFailure)
	}
	if hso["additionalContext"] != "retry with smaller input" {
		t.Fatalf("additionalContext = %v", hso["additionalContext"])
	}
}
