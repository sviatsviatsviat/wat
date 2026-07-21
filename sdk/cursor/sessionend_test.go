package cursor

import (
	"testing"
)

func TestDecode_SessionEnd(t *testing.T) {
	e := mustDecode[SessionEnd](t, `{"hook_event_name":"sessionEnd","conversation_id":"c1","reason":"complete","is_background_agent":false}`)
	if e.Reason != "complete" {
		t.Fatalf("event=%+v", e)
	}
}
