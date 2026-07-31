package stopevent

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestEncode_StopBlock(t *testing.T) {
	out, _, err := NewResults(event.Stop).FollowUp("run the tests").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_Stop(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"Stop","stop_hook_active":false,"last_assistant_message":"bye"}`, event.Stop)
	if ev.StopHookActive {
		t.Fatal("StopHookActive = true, want false")
	}
	if ev.LastAssistantMessage != "bye" {
		t.Fatalf("LastAssistantMessage = %q, want bye", ev.LastAssistantMessage)
	}
	if ev.BackgroundTasks != nil {
		t.Fatalf("BackgroundTasks = %#v, want nil", ev.BackgroundTasks)
	}
	if ev.SessionCrons != nil {
		t.Fatalf("SessionCrons = %#v, want nil", ev.SessionCrons)
	}
}

func TestDecode_StopBackgroundTasksAndCrons(t *testing.T) {
	raw := `{
		"session_id":"s",
		"hook_event_name":"Stop",
		"stop_hook_active":true,
		"last_assistant_message":"I've completed the refactoring.",
		"background_tasks":[
			{
				"id":"task-001",
				"type":"shell",
				"status":"running",
				"description":"tail logs",
				"command":"tail -f /var/log/syslog"
			},
			{
				"id":"task-002",
				"type":"subagent",
				"status":"running",
				"description":"explore",
				"agent_type":"Explore"
			},
			{
				"id":"task-003",
				"type":"MCP task",
				"status":"running",
				"server":"github",
				"tool":"list_issues"
			},
			{
				"id":"task-004",
				"type":"workflow",
				"status":"running",
				"name":"ci"
			}
		],
		"session_crons":[
			{
				"id":"cron-001",
				"schedule":"0 9 * * 1-5",
				"recurring":true,
				"prompt":"check the build"
			},
			{
				"id":"cron-002",
				"schedule":"0 12 30 7 *",
				"recurring":false,
				"prompt":"one-shot wakeup"
			}
		]
	}`
	ev := mustDecode[Event](t, raw, event.Stop)
	if !ev.StopHookActive {
		t.Fatal("StopHookActive = false, want true")
	}
	if len(ev.BackgroundTasks) != 4 {
		t.Fatalf("len(BackgroundTasks) = %d, want 4", len(ev.BackgroundTasks))
	}
	shell := ev.BackgroundTasks[0]
	if shell.ID != "task-001" || shell.Type != "shell" || shell.Status != "running" ||
		shell.Description != "tail logs" || shell.Command != "tail -f /var/log/syslog" {
		t.Fatalf("shell task = %#v", shell)
	}
	sub := ev.BackgroundTasks[1]
	if sub.AgentType != "Explore" || sub.Command != "" {
		t.Fatalf("subagent task = %#v", sub)
	}
	mcp := ev.BackgroundTasks[2]
	if mcp.Server != "github" || mcp.Tool != "list_issues" {
		t.Fatalf("mcp task = %#v", mcp)
	}
	wf := ev.BackgroundTasks[3]
	if wf.Name != "ci" {
		t.Fatalf("workflow task = %#v", wf)
	}
	if len(ev.SessionCrons) != 2 {
		t.Fatalf("len(SessionCrons) = %d, want 2", len(ev.SessionCrons))
	}
	cron := ev.SessionCrons[0]
	if cron.ID != "cron-001" || cron.Schedule != "0 9 * * 1-5" || !cron.Recurring || cron.Prompt != "check the build" {
		t.Fatalf("session cron = %#v", cron)
	}
	oneshot := ev.SessionCrons[1]
	if oneshot.Recurring || oneshot.Prompt != "one-shot wakeup" {
		t.Fatalf("one-shot cron = %#v", oneshot)
	}
}

func TestDecode_StopEmptyBackgroundArrays(t *testing.T) {
	ev := mustDecode[Event](t, `{
		"session_id":"s",
		"hook_event_name":"Stop",
		"stop_hook_active":false,
		"last_assistant_message":"done",
		"background_tasks":[],
		"session_crons":[]
	}`, event.Stop)
	if ev.BackgroundTasks == nil || len(ev.BackgroundTasks) != 0 {
		t.Fatalf("BackgroundTasks = %#v, want empty non-nil", ev.BackgroundTasks)
	}
	if ev.SessionCrons == nil || len(ev.SessionCrons) != 0 {
		t.Fatalf("SessionCrons = %#v, want empty non-nil", ev.SessionCrons)
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
