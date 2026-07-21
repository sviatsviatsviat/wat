package copilot

import "testing"

func TestDecode_Notification(t *testing.T) {
	raw := `{
  "hook_event_name": "Notification",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "message": "shell done",
  "title": "Shell completed",
  "notification_type": "shell_completed"
}`
	ev, err := codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	note, ok := ev.(Notification)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if note.EventName() != EventNotification {
		t.Fatalf("EventName=%q", note.EventName())
	}
	if note.NotificationType != "shell_completed" || note.Message != "shell done" {
		t.Fatalf("Notification=%+v", note)
	}
}
