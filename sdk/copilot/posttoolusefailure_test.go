package copilot

import (
	"testing"

)

func TestEncode_PostToolFailureContext(t *testing.T) {
	out, code, err := postToolFailureResults{}.Context("retry with smaller input").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != WarnExit {
		t.Fatalf("code=%d, want %d", code, WarnExit)
	}
	if string(out) != "retry with smaller input" {
		t.Fatalf("stdout=%q", out)
	}
}

func TestDecode_PostToolUseFailure(t *testing.T) {
	e := mustDecode[PostToolUseFailure](t, `{"hook_event_name":"PostToolUseFailure","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","tool_name":"bash","tool_input":{},"error":"timeout"}`, EventPostToolUseFailure)
	if e.ErrorMessage() != "timeout" {
		t.Fatalf("PostToolUseFailure=%+v", e)
	}
}
