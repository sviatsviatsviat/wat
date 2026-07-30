package worktreecreate

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func TestDecode_WorktreeCreate(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"WorktreeCreate","name":"feature-auth"}`, event.WorktreeCreate)
	if ev.Name != "feature-auth" {
		t.Fatalf("Name = %q", ev.Name)
	}
}

func TestEncode_Path_plainStdout(t *testing.T) {
	out, code, err := results{}.Path("/tmp/wt").Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
	}
	if string(out) != "/tmp/wt" {
		t.Fatalf("stdout = %q, want plain path", out)
	}
	if json.Valid(out) {
		t.Fatalf("command-hook path must not be JSON, got %s", out)
	}
}

func TestEncode_Path_ignoresCommonJSON(t *testing.T) {
	out, code, err := results{}.Path("/tmp/wt").WithSystemMessage("note").WithContinue(true).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
	}
	if string(out) != "/tmp/wt" {
		t.Fatalf("stdout = %q, want plain path only", out)
	}
	if bytes.Contains(out, []byte("systemMessage")) || bytes.Contains(out, []byte("hookSpecificOutput")) {
		t.Fatalf("Common/JSON fields must not appear on command-hook stdout: %s", out)
	}
}

func TestEncode_emptyPath(t *testing.T) {
	out, code, err := (output{}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
	}
	if out != nil {
		t.Fatalf("empty path encode = %q, want nil", out)
	}
	if !(output{}).IsZero() {
		t.Fatal("empty worktree create output should be zero")
	}
}

func TestMerge_Path_lastWins(t *testing.T) {
	a := results{}.Path("/tmp/a")
	b := results{}.Path("/tmp/b")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] == "" {
		t.Fatalf("warnings = %v, want overwrite warning", warnings)
	}
	out, code, err := merged.(Output).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d, want %d", code, event.SuccessExit)
	}
	if string(out) != "/tmp/b" {
		t.Fatalf("merged path = %q, want /tmp/b", out)
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
