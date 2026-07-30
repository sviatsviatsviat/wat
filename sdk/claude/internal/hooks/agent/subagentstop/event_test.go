package subagentstop

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_SubagentStop(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"SubagentStop","agent_id":"a1","agent_type":"worker","stop_hook_active":true,"last_assistant_message":"done"}`, event.SubagentStop)
	if !ev.StopHookActive {
		t.Fatal("StopHookActive = false, want true")
	}
	if ev.LastAssistantMessage != "done" {
		t.Fatalf("LastAssistantMessage = %q, want done", ev.LastAssistantMessage)
	}
	if ev.AgentID != "a1" || ev.AgentType != "worker" {
		t.Fatalf("agent identity = %q/%q", ev.AgentID, ev.AgentType)
	}
	if ev.BackgroundTasks != nil || ev.SessionCrons != nil {
		t.Fatalf("background fields = %#v %#v, want nil", ev.BackgroundTasks, ev.SessionCrons)
	}
}

func TestDecode_SubagentStopBackgroundTasksAndCrons(t *testing.T) {
	raw := `{
		"session_id":"s",
		"hook_event_name":"SubagentStop",
		"agent_id":"a1",
		"agent_type":"worker",
		"agent_transcript_path":"/tmp/sub.jsonl",
		"stop_hook_active":false,
		"last_assistant_message":"subagent done",
		"background_tasks":[
			{
				"id":"task-parent",
				"type":"monitor",
				"status":"running",
				"description":"watch",
				"server":"sentry",
				"tool":"tail"
			}
		],
		"session_crons":[
			{
				"id":"cron-parent",
				"schedule":"*/5 * * * *",
				"recurring":true,
				"prompt":"poll"
			}
		]
	}`
	ev := mustDecode[Event](t, raw, event.SubagentStop)
	if ev.AgentTranscriptPath != "/tmp/sub.jsonl" {
		t.Fatalf("AgentTranscriptPath = %q", ev.AgentTranscriptPath)
	}
	if len(ev.BackgroundTasks) != 1 {
		t.Fatalf("len(BackgroundTasks) = %d, want 1", len(ev.BackgroundTasks))
	}
	task := ev.BackgroundTasks[0]
	if task.ID != "task-parent" || task.Type != "monitor" || task.Server != "sentry" || task.Tool != "tail" {
		t.Fatalf("background task = %#v", task)
	}
	if len(ev.SessionCrons) != 1 {
		t.Fatalf("len(SessionCrons) = %d, want 1", len(ev.SessionCrons))
	}
	cron := ev.SessionCrons[0]
	if cron.ID != "cron-parent" || cron.Schedule != "*/5 * * * *" || !cron.Recurring || cron.Prompt != "poll" {
		t.Fatalf("session cron = %#v", cron)
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
