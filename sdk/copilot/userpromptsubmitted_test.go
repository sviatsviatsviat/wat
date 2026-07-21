package copilot

import (
	"testing"
)

func TestDecode_UserPromptSubmitted(t *testing.T) {
	e := mustDecode[UserPromptSubmitted](t, `{"hook_event_name":"UserPromptSubmit","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","prompt":"hello"}`, EventUserPromptSubmitted)
	if e.Prompt != "hello" {
		t.Fatalf("UserPromptSubmitted=%+v", e)
	}
}
