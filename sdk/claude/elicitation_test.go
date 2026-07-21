package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_Elicitation(t *testing.T) {
	mustDecode[Elicitation](t, `{"session_id":"s","hook_event_name":"Elicitation","server_name":"srv","message":"confirm?"}`, EventElicitation)

	out, code, err := elicitationResults{}.Accept().WithContent(map[string]any{"answer": "yes"}).Encode()
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
	if hso["hookEventName"] != EventElicitation {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventElicitation)
	}
	if hso["action"] != "accept" {
		t.Fatalf("action = %v", hso["action"])
	}
	content, ok := hso["content"].(map[string]any)
	if !ok || content["answer"] != "yes" {
		t.Fatalf("content = %v", hso["content"])
	}

	// ElicitationResults has no Noop(); empty output is the zero-state contract.
	if !(elicitationOutput{}).IsZero() {
		t.Fatal("empty elicitation output should be zero")
	}
}
