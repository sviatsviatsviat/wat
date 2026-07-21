package copilot

import (
	"strings"
	"testing"
)

func TestEncode_SessionStartContext(t *testing.T) {
	out, code, err := codec.Encode(sessionStartResults{}.Context("project uses go test ./..."))
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
