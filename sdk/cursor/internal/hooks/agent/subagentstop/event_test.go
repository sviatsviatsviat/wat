package subagentstop

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

const cursorSubagentStop = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "default",
  "hook_event_name": "subagentStop",
  "cursor_version": "3.13.10",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "cwd": "/w",
  "subagent_id": "sa1",
  "subagent_type": "generalPurpose",
  "status": "completed",
  "task": "Explore the authentication flow",
  "description": "Exploring auth flow",
  "summary": "Auth uses JWT middleware",
  "duration_ms": 45000,
  "message_count": 12,
  "tool_call_count": 8,
  "loop_count": 0,
  "modified_files": ["src/auth.ts"],
  "agent_transcript_path": "/path/to/subagent/transcript.txt"
}`

func TestDecode_SubagentStop(t *testing.T) {
	ev := mustDecode[Event](t, cursorSubagentStop)

	if ev.SubagentID != "sa1" {
		t.Errorf("SubagentID = %q, want sa1", ev.SubagentID)
	}
	if ev.SubagentType != "generalPurpose" {
		t.Errorf("SubagentType = %q, want generalPurpose", ev.SubagentType)
	}
	if ev.Status != "completed" {
		t.Errorf("Status = %q, want completed", ev.Status)
	}
	if ev.Task != "Explore the authentication flow" {
		t.Errorf("Task = %q", ev.Task)
	}
	if ev.Description != "Exploring auth flow" {
		t.Errorf("Description = %q, want Exploring auth flow", ev.Description)
	}
	if ev.Summary != "Auth uses JWT middleware" {
		t.Errorf("Summary = %q", ev.Summary)
	}
	if ev.DurationMs != 45000 {
		t.Errorf("DurationMs = %d, want 45000", ev.DurationMs)
	}
	if ev.MessageCount != 12 {
		t.Errorf("MessageCount = %d, want 12", ev.MessageCount)
	}
	if ev.ToolCallCount != 8 {
		t.Errorf("ToolCallCount = %d, want 8", ev.ToolCallCount)
	}
	if ev.LoopCount != 0 {
		t.Errorf("LoopCount = %d, want 0", ev.LoopCount)
	}
	wantFiles := []string{"src/auth.ts"}
	if !reflect.DeepEqual(ev.ModifiedFiles, wantFiles) {
		t.Errorf("ModifiedFiles = %#v, want %#v", ev.ModifiedFiles, wantFiles)
	}
	if ev.AgentTranscriptPath == nil || *ev.AgentTranscriptPath != "/path/to/subagent/transcript.txt" {
		t.Errorf("AgentTranscriptPath = %v", ev.AgentTranscriptPath)
	}
}

func TestDecode_SubagentStop_nullTranscript(t *testing.T) {
	ev := mustDecode[Event](t, `{"hook_event_name":"subagentStop","conversation_id":"c1","subagent_type":"explore","status":"aborted","task":"t","description":"d","summary":"","duration_ms":1,"message_count":0,"tool_call_count":0,"loop_count":1,"modified_files":[],"agent_transcript_path":null}`)
	if ev.AgentTranscriptPath != nil {
		t.Fatalf("AgentTranscriptPath = %v, want nil", ev.AgentTranscriptPath)
	}
	if ev.Status != "aborted" {
		t.Fatalf("Status = %q, want aborted", ev.Status)
	}
	if ev.Description != "d" {
		t.Fatalf("Description = %q, want d", ev.Description)
	}
	if len(ev.ModifiedFiles) != 0 {
		t.Fatalf("ModifiedFiles = %#v, want empty", ev.ModifiedFiles)
	}
}

func TestEncode_SubagentStopFollowUp(t *testing.T) {
	out, code, err := stopevent.NewResults().FollowUp("continue with next step").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if !strings.Contains(string(out), `"followup_message":"continue with next step"`) {
		t.Fatalf("bad subagentStop output: %s", out)
	}
}

func init() {
	register(testCodec)
}
