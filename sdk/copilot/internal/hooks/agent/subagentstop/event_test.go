package subagentstop

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func init() {
	register(testCodec)
}

const copilotSubagentStop = `{
  "hook_event_name": "SubagentStop",
  "session_id": "s",
  "timestamp": "2026-01-01T00:00:00Z",
  "cwd": "/w",
  "transcript_path": "/t",
  "agent_name": "task",
  "agent_display_name": "Task",
  "stop_reason": "end_turn",
  "last_assistant_message": "full final subagent response"
}`

func TestDecode_SubagentStop(t *testing.T) {
	ev, err := testCodec.Decode([]byte(copilotSubagentStop))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok {
		t.Fatalf("got %T", ev)
	}
	if e.Name() != "task" || e.DisplayName() != "Task" {
		t.Fatalf("identity Name=%q DisplayName=%q", e.Name(), e.DisplayName())
	}
	if e.Reason() != "end_turn" || e.EventName() != event.SubagentStop {
		t.Fatalf("Reason=%q EventName=%q", e.Reason(), e.EventName())
	}
	if e.LastAssistantMessage != "full final subagent response" {
		t.Fatalf("LastAssistantMessage=%q", e.LastAssistantMessage)
	}
	if e.HookEventName != "SubagentStop" {
		t.Fatalf("HookEventName=%q", e.HookEventName)
	}
}

func TestDecode_SubagentStop_withoutMessage(t *testing.T) {
	ev, err := testCodec.Decode([]byte(`{"hook_event_name":"SubagentStop","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","transcript_path":"/t","agent_name":"task","stop_reason":"end_turn"}`))
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ev.(Event)
	if !ok || e.LastAssistantMessage != "" {
		t.Fatalf("SubagentStop=%+v", ev)
	}
}

func TestEncode_SubagentStopFollowUp(t *testing.T) {
	out, code, err := NewResults().FollowUp("continue").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "continue") {
		t.Fatalf("bad output: %s", out)
	}
	if strings.Contains(string(out), "modifiedResponse") {
		t.Fatalf("block output should omit modifiedResponse: %s", out)
	}
}

func TestEncode_SubagentStopModifiedResponse(t *testing.T) {
	out, code, err := NewResults().ModifiedResponse("redacted for parent").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got["modifiedResponse"] != "redacted for parent" {
		t.Fatalf("got %#v", got)
	}
	if _, ok := got["decision"]; ok {
		t.Fatalf("unexpected decision in %#v", got)
	}
}

func TestEncode_SubagentStopBlockWinsOverModifiedResponse(t *testing.T) {
	out, code, err := NewResults().FollowUp("keep going").WithModifiedResponse("ignored").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "keep going") {
		t.Fatalf("bad block output: %s", out)
	}
	if strings.Contains(string(out), "modifiedResponse") || strings.Contains(string(out), "ignored") {
		t.Fatalf("block must discard modifiedResponse: %s", out)
	}
}

func TestMerge_SubagentStopTakesLast(t *testing.T) {
	a := NewResults().ModifiedResponse("first")
	b := NewResults().ModifiedResponse("second")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "modifiedResponse") {
		t.Fatalf("warnings=%v", warnings)
	}
	out, code, err := merged.Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"modifiedResponse":"second"`) {
		t.Fatalf("merged=%s", out)
	}
	if merged.Stop() {
		t.Fatal("modifiedResponse alone must not Stop the chain")
	}
}

func TestMerge_SubagentStopFollowUpStops(t *testing.T) {
	a := NewResults().ModifiedResponse("rewrite")
	b := NewResults().FollowUp("continue")
	merged, _, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if !merged.Stop() {
		t.Fatal("FollowUp must Stop the chain")
	}
	out, _, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || strings.Contains(string(out), "modifiedResponse") {
		t.Fatalf("encode=%s", out)
	}
}

func TestEncode_SubagentStopNoop(t *testing.T) {
	out, code, err := NewResults().Noop().Encode()
	if err != nil || code != 0 || out != nil {
		t.Fatalf("Noop Encode = %q %d %v", out, code, err)
	}
}
