package taskcompleted

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_TaskCompleted(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"TaskCompleted",
		"task_id":"task-001",
		"task_subject":"Implement user authentication",
		"task_description":"Add login and signup endpoints",
		"teammate_name":"implementer",
		"team_name":"session-a1b2c3d4"
	}`, event.TaskCompleted)
	if ev.TaskID != "task-001" {
		t.Fatalf("TaskID = %q", ev.TaskID)
	}
	if ev.TaskSubject != "Implement user authentication" {
		t.Fatalf("TaskSubject = %q", ev.TaskSubject)
	}
	if ev.TaskDescription != "Add login and signup endpoints" {
		t.Fatalf("TaskDescription = %q", ev.TaskDescription)
	}
	if ev.TeammateName != "implementer" {
		t.Fatalf("TeammateName = %q", ev.TeammateName)
	}
	if ev.TeamName != "session-a1b2c3d4" {
		t.Fatalf("TeamName = %q", ev.TeamName)
	}
}

func init() {
	register(testCodec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != wantName {
		t.Fatalf("EventName() = %q, want %q", ev.EventName(), wantName)
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}
