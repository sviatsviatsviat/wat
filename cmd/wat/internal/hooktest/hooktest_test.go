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

func TestResolveFixture(t *testing.T) {
	tests := []struct {
		name        string
		agent       string
		payload     string
		wantDialect string
		wantEvent   string
		wantErr     string
	}{
		{
			name:    "copilot_requires_hook_event_name",
			agent:   "copilot",
			payload: `{"session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`,
			wantErr: "event name",
		},
		{
			name:        "claude_pre_tool_use",
			agent:       "claude",
			payload:     `{"hook_event_name":"PreToolUse","session_id":"s1"}`,
			wantDialect: sdkclaude.Dialect,
			wantEvent:   "PreToolUse",
		},
		{
			name:    "unknown_agent",
			agent:   "nosuch",
			payload: `{"hook_event_name":"PreToolUse"}`,
			wantErr: "unknown dialect",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ResolveFixture(tt.agent, []byte(tt.payload))
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				errText := err.Error()
				ok := strings.Contains(errText, tt.wantErr)
				if tt.wantErr == "event name" {
					ok = ok || strings.Contains(errText, "hook_event_name")
				}
				if !ok {
					t.Fatalf("error = %q, want substring %q", errText, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if info.Dialect != tt.wantDialect {
				t.Fatalf("dialect = %q, want %q", info.Dialect, tt.wantDialect)
			}
			if info.Event != tt.wantEvent {
				t.Fatalf("event = %q, want %q", info.Event, tt.wantEvent)
			}
		})
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
