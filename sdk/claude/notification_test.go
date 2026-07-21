package claude

import (
	"encoding/json"
	"testing"
)

func TestDecode_Notification(t *testing.T) {
	mustDecode[Notification](t, `{"session_id":"s","hook_event_name":"Notification","notification_type":"idle_prompt","message":"waiting"}`, EventNotification)

	out, code, err := notificationResults{}.Context("idle note").Encode()
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
	if hso["hookEventName"] != EventNotification {
		t.Fatalf("hookEventName = %v, want %q", hso["hookEventName"], EventNotification)
	}
	if hso["additionalContext"] != "idle note" {
		t.Fatalf("additionalContext = %v", hso["additionalContext"])
	}
}
