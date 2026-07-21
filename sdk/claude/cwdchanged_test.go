package claude

import (
	"testing"
)

func TestDecode_CwdChanged(t *testing.T) {
	ev := mustDecode[CwdChanged](t, `{"session_id":"s","hook_event_name":"CwdChanged","new_cwd":"/new","old_cwd":"/old"}`, EventCwdChanged)
	if ev.NewCwd != "/new" || ev.OldCwd != "/old" {
		t.Fatalf("cwd fields = %+v", ev)
	}
}
