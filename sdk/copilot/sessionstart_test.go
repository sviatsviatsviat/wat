package copilot

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := sessionStartResults{}.Context("project uses go test ./...").Encode()
	if err != nil || code != 0 {
		t.Fatal(err, code)
	}
	if !strings.Contains(string(out), `"additional_context"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestDecode_SessionStart(t *testing.T) {
	e := mustDecode[SessionStart](t, `{"hook_event_name":"SessionStart","session_id":"s","timestamp":"2026-01-01T00:00:00Z","cwd":"/w","source":"new","initial_prompt":"go"}`, EventSessionStart)
	if e.Source != "new" || e.InitialPrompt() != "go" {
		t.Fatalf("SessionStart=%+v", e)
	}
}

func TestMerge_SessionStart_contextJoins(t *testing.T) {
	a := sessionStartResults{}.Context("one")
	b := sessionStartResults{}.Context("two")
	merged, warnings, err := a.Merge(b.(run.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(sessionStartOutput)
	if out.additionalContext != "one\n\ntwo" {
		t.Fatalf("context = %q", out.additionalContext)
	}
	if merged.Stop() {
		t.Fatal("context should not stop")
	}
}
