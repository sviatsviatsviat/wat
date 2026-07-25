package subagentstart

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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

func TestDecode_SubagentStart(t *testing.T) {
	ev := mustDecode[Event](t, `{"hook_event_name":"subagentStart","conversation_id":"c1","subagent_id":"sa1","subagent_type":"explore","task":"find files","parent_conversation_id":"c0","tool_call_id":"tc1","subagent_model":"claude-4.5-sonnet","is_parallel_worker":true,"git_branch":"main"}`)

	if ev.SubagentID != "sa1" {
		t.Errorf("SubagentID = %q, want sa1", ev.SubagentID)
	}
	if ev.SubagentType != "explore" {
		t.Errorf("SubagentType = %q, want explore", ev.SubagentType)
	}
	if ev.Task != "find files" {
		t.Errorf("Task = %q, want %q", ev.Task, "find files")
	}
	if ev.ParentConversationID != "c0" {
		t.Errorf("ParentConversationID = %q, want c0", ev.ParentConversationID)
	}
	if ev.ToolCallID != "tc1" {
		t.Errorf("ToolCallID = %q, want tc1", ev.ToolCallID)
	}
	if ev.SubagentModel != "claude-4.5-sonnet" {
		t.Errorf("SubagentModel = %q, want claude-4.5-sonnet", ev.SubagentModel)
	}
	if !ev.IsParallelWorker {
		t.Error("IsParallelWorker = false, want true")
	}
	if ev.GitBranch != "main" {
		t.Errorf("GitBranch = %q, want main", ev.GitBranch)
	}
}

func TestEncode_SubagentStartDeny_userMessageExitZero(t *testing.T) {
	out, code, err := results{}.Deny("pinned model blocked").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	got := string(out)
	if !strings.Contains(got, `"permission":"deny"`) {
		t.Fatalf("missing permission deny: %s", got)
	}
	if !strings.Contains(got, `"user_message":"pinned model blocked"`) {
		t.Fatalf("missing user_message: %s", got)
	}
	if strings.Contains(got, "agent_message") {
		t.Fatalf("subagentStart deny must not emit agent_message: %s", got)
	}
}

func init() {
	register(testCodec)
}
