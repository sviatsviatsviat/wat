package cursor

import (
	"testing"
)

func TestDecode_PostToolUse(t *testing.T) {
	e := mustDecode[PostToolUse](t, `{"hook_event_name":"postToolUse","conversation_id":"c1","tool_name":"Read","tool_output":"contents","duration":100}`)
	if e.ToolOutput != "contents" || e.DurationMillis() != 100 {
		t.Fatalf("event=%+v", e)
	}
}
