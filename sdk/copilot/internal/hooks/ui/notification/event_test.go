package notification

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

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
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	note, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if note.EventName() != event.Notification {
		t.Fatalf("EventName=%q", note.EventName())
	}
	if note.NotificationType != "shell_completed" || note.Message != "shell done" {
		t.Fatalf("Notification=%+v", note)
	}
}

func TestEncode_NotificationContext(t *testing.T) {
	out, code, err := results{}.Context("extra").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additional_context":"extra"`) {
		t.Fatalf("bad output: %s", out)
	}
}
