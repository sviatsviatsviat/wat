package hooktest

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func TestResolveFixture_copilotRequiresHookEventName(t *testing.T) {
	payload := []byte(`{"session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`)

	_, err := ResolveFixture("copilot", payload)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "hook_event_name") && !strings.Contains(err.Error(), "event name") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestResolveFixture_claudePreToolUse(t *testing.T) {
	payload := []byte(`{"hook_event_name":"PreToolUse","session_id":"s1"}`)
	info, err := ResolveFixture("claude", payload)
	if err != nil {
		t.Fatal(err)
	}
	if info.Dialect != sdkclaude.Dialect {
		t.Fatalf("dialect = %q, want %q", info.Dialect, sdkclaude.Dialect)
	}
	if info.Event != "PreToolUse" {
		t.Fatalf("event = %q, want PreToolUse", info.Event)
	}
}

func TestResolveFixture_unknownAgent(t *testing.T) {
	payload := []byte(`{"hook_event_name":"PreToolUse"}`)
	_, err := ResolveFixture("nosuch", payload)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown dialect") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestWriteReport_fixtureSummary(t *testing.T) {
	var buf bytes.Buffer
	WriteReport(&buf, FixtureInfo{Dialect: sdkclaude.Dialect, Event: "PreToolUse"}, []byte(`{"permissionDecision":"deny","reason":"blocked"}`), nil, 0, false)

	out := buf.String()
	for _, want := range []string{"agent: claude", "event: PreToolUse", "decision: deny", "exit:   0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %q:\n%s", want, out)
		}
	}
}

func TestRun_emptyFixture(t *testing.T) {
	project := t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var errBuf bytes.Buffer
	deps := DefaultDeps(io.Discard)
	deps.Getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return project
		}
		return ""
	}
	deps.ReadFixture = func(string, io.Reader) ([]byte, error) {
		return nil, nil
	}

	code := Run(Config{Fixture: "-"}, "v0.0.0-test", deps, strings.NewReader(""), &errBuf)
	if code != ExitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, ExitRuntimeFailure)
	}
	if !strings.Contains(errBuf.String(), "empty fixture") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}
