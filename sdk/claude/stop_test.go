package claude

import (
	"strings"
	"testing"
)

func TestEncode_StopBlock(t *testing.T) {
	out, _, err := stopResults{eventName: EventStop}.FollowUp("run the tests").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_Stop(t *testing.T) {
	mustDecode[Stop](t, `{"session_id":"s","hook_event_name":"Stop","stop_hook_active":false,"last_assistant_message":"bye"}`, EventStop)
}
