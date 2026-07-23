package agentstop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

const copilotVSCodeStop = `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "transcript_path": "/tmp/t",
  "stop_reason": "end_turn"
}`

func TestDecode_VSCodeStop(t *testing.T) {
	ev, err := testCodec.Decode([]byte(copilotVSCodeStop))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if stop.EventName() != event.AgentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if stop.HookEventName != "Stop" {
		t.Fatalf("HookEventName=%q", stop.HookEventName)
	}
	if stop.Reason() != "end_turn" {
		t.Fatalf("Reason=%q", stop.Reason())
	}
	if stop.TranscriptPath != "/tmp/t" {
		t.Fatalf("TranscriptPath=%q", stop.TranscriptPath)
	}
}

func TestEncode_StopFollowUp(t *testing.T) {
	out, code, err := results{}.FollowUp("run the tests").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_VSCodeStopWithAgentScope(t *testing.T) {
	raw := `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "agent_name": "task",
  "stop_reason": "end_turn"
}`
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	stop, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T, want Event", ev)
	}
	if stop.EventName() != event.AgentStop {
		t.Fatalf("EventName=%q", stop.EventName())
	}
	if !stop.IsSubagent() || stop.Name() != "task" {
		t.Fatalf("IsSubagent/Name = %v %q", stop.IsSubagent(), stop.Name())
	}
	if stop.HookEventName != "Stop" {
		t.Fatalf("HookEventName=%q", stop.HookEventName)
	}
}

func TestDecode_AgentStop(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"Stop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","stop_reason":"end_turn"}`))
	if err != nil {
		t.Fatal(err)
	}
	e := ev.(Event)
	if e.Reason() != "end_turn" || e.IsSubagent() {
		t.Fatalf("AgentStop=%+v", e)
	}
}

func TestDecode_AgentStopWithAgentScope(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"Stop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`))
	if err != nil {
		t.Fatal(err)
	}
	e := ev.(Event)
	if !e.IsSubagent() || e.Name() != "task" {
		t.Fatalf("AgentStop=%+v", e)
	}
}

func init() {
	register(testCodec)
}
