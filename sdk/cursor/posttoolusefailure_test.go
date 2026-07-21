package cursor

import (
	"testing"
)

func TestDecode_PostToolUseFailure(t *testing.T) {
	e := mustDecode[PostToolUseFailure](t, `{"hook_event_name":"postToolUseFailure","conversation_id":"c1","tool_name":"Shell","error_message":"timeout","failure_type":"timeout","duration_ms":50}`)
	if e.ErrorMessage != "timeout" || e.FailureType != "timeout" {
		t.Fatalf("event=%+v", e)
	}
}
