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

func TestSidecarExpectPath(t *testing.T) {
	tests := []struct {
		fixture string
		want    string
	}{
		{fixture: "", want: ""},
		{fixture: "-", want: ""},
		{fixture: "foo.json", want: "foo.expect.json"},
		{fixture: "dir/subagent_start.json", want: "dir/subagent_start.expect.json"},
		{fixture: "payload", want: "payload.expect.json"},
	}
	for _, tt := range tests {
		if got := SidecarExpectPath(tt.fixture); got != tt.want {
			t.Fatalf("SidecarExpectPath(%q) = %q, want %q", tt.fixture, got, tt.want)
		}
	}
}

func TestResolveExpectPath(t *testing.T) {
	dir := t.TempDir()
	fixture := filepath.Join(dir, "case.json")
	sidecar := filepath.Join(dir, "case.expect.json")
	explicit := filepath.Join(dir, "custom.expect.json")
	if err := os.WriteFile(sidecar, []byte(`{"exit":0}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(explicit, []byte(`{"exit":2}`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("explicit wins", func(t *testing.T) {
		got, err := ResolveExpectPath(fixture, explicit, os.Stat)
		if err != nil {
			t.Fatal(err)
		}
		if got != explicit {
			t.Fatalf("got %q, want %q", got, explicit)
		}
	})
	t.Run("sidecar when present", func(t *testing.T) {
		got, err := ResolveExpectPath(fixture, "", os.Stat)
		if err != nil {
			t.Fatal(err)
		}
		if got != sidecar {
			t.Fatalf("got %q, want %q", got, sidecar)
		}
	})
	t.Run("missing sidecar is optional", func(t *testing.T) {
		got, err := ResolveExpectPath(filepath.Join(dir, "missing.json"), "", os.Stat)
		if err != nil {
			t.Fatal(err)
		}
		if got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})
	t.Run("explicit missing errors", func(t *testing.T) {
		_, err := ResolveExpectPath(fixture, filepath.Join(dir, "nope.json"), os.Stat)
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestCheckExpect(t *testing.T) {
	denyExit := 2
	empty := true
	tests := []struct {
		name    string
		exp     Expect
		stdout  string
		exit    int
		wantLen int
	}{
		{
			name: "pass deny",
			exp: Expect{
				Exit:           &denyExit,
				Decision:       "deny",
				StdoutContains: []string{"blocked"},
			},
			stdout:  `{"permission":"deny","user_message":"blocked"}`,
			exit:    2,
			wantLen: 0,
		},
		{
			name:    "fail exit",
			exp:     Expect{Exit: &denyExit},
			stdout:  `{}`,
			exit:    0,
			wantLen: 1,
		},
		{
			name:    "fail empty",
			exp:     Expect{StdoutEmpty: &empty},
			stdout:  `{"permission":"allow"}`,
			exit:    0,
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fails := CheckExpect(tt.exp, "cursor", []byte(tt.stdout), tt.exit)
			if len(fails) != tt.wantLen {
				t.Fatalf("fails = %#v, want len %d", fails, tt.wantLen)
			}
		})
	}
}

func TestLoadExpect_rejectsUnknownField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.expect.json")
	if err := os.WriteFile(path, []byte(`{"exitt":0}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadExpect(path, os.ReadFile)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "decode expect") {
		t.Fatalf("error = %q", err)
	}
}
