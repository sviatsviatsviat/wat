package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hooktest"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/initproj"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

func TestRunTest_missingFixtureUsage(t *testing.T) {
	prevStderr := stderr
	var errBuf bytes.Buffer
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	code := run([]string{"test"})
	if code != exitUsage {
		t.Fatalf("exit = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(errBuf.String(), "--fixture is required") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestRunTest_emptyFixture(t *testing.T) {
	prevStderr := stderr
	var errBuf bytes.Buffer
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	project := t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	deps := hooktest.DefaultDeps(io.Discard)
	deps.Getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return project
		}
		return ""
	}
	deps.ReadFixture = func(string, io.Reader) ([]byte, error) {
		return nil, nil
	}

	code := hooktest.Run(hooktest.Config{Fixture: "-"}, watModuleVersionFn(), deps, os.Stdin, stderr)
	if code != hooktest.ExitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, hooktest.ExitRuntimeFailure)
	}
	if !strings.Contains(errBuf.String(), "empty fixture") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestResolveFixture_copilotRequiresHookEventName(t *testing.T) {
	payload := []byte(`{"session_id":"s1","timestamp":"2026-07-12T10:00:00Z","cwd":"/w"}`)

	_, err := hooktest.ResolveFixture("copilot", payload)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "hook_event_name") && !strings.Contains(err.Error(), "event name") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestWriteTestReport_fixtureSummary(t *testing.T) {
	var buf bytes.Buffer
	hooktest.WriteReport(&buf, hooktest.FixtureInfo{Dialect: sdkclaude.Dialect, Event: "PreToolUse"}, []byte(`{"permissionDecision":"deny","reason":"blocked"}`), nil, 0, false)

	out := buf.String()
	for _, want := range []string{"agent: claude", "event: PreToolUse", "decision: deny", "exit:   0"} {
		if !strings.Contains(out, want) {
			t.Fatalf("report missing %q:\n%s", want, out)
		}
	}
}

func TestRunTest_preToolDenyIntegration(t *testing.T) {
	project := setupTestHookProject(t)
	fixtures := testModuleRoot(t)

	tests := []struct {
		name       string
		agent      string
		event      string
		fixture    string
		wantExit   int
		wantOutput []string
	}{
		{
			name:     "claude",
			agent:    "claude",
			fixture:  filepath.Join(fixtures, "testdata", "fixtures", "claude", "pre_tool_force_push.json"),
			wantExit: exitOK,
			wantOutput: []string{
				"event: PreToolUse",
				"deny",
				"force pushes are not allowed",
			},
		},
		{
			name:     "copilot",
			agent:    "copilot",
			fixture:  filepath.Join(fixtures, "testdata", "fixtures", "copilot", "pre_tool_force_push.json"),
			wantExit: exitOK,
			wantOutput: []string{
				`"permission_decision":"deny"`,
				"force pushes are not allowed",
			},
		},
		{
			name:     "cursor",
			agent:    "cursor",
			fixture:  filepath.Join(fixtures, "testdata", "fixtures", "cursor", "before_shell_force_push.json"),
			wantExit: cursor.PermissionDenyExit,
			wantOutput: []string{
				`"permission":"deny"`,
				"force pushes are not allowed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var report bytes.Buffer
			deps := hooktest.DefaultDeps(&report)
			deps.Getenv = func(key string) string {
				if key == "WAT_PROJECT_DIR" {
					return project
				}
				return os.Getenv(key)
			}

			code := hooktest.Run(hooktest.Config{
				Agent:   tt.agent,
				Fixture: tt.fixture,
			}, watModuleVersionFn(), deps, os.Stdin, stderr)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d\nreport:\n%s", code, tt.wantExit, report.String())
			}
			out := report.String()
			for _, want := range tt.wantOutput {
				if !strings.Contains(out, want) {
					t.Fatalf("report missing %q:\n%s", want, out)
				}
			}
		})
	}
}

func setupTestHookProject(t *testing.T) string {
	t.Helper()

	prevVersionFn := watModuleVersionFn
	watModuleVersionFn = func() string { return "v0.0.0-test-000000000000" }
	t.Cleanup(func() { watModuleVersionFn = prevVersionFn })

	dir := t.TempDir()
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	if err := initproj.Init(dir, false, watModuleVersionFn(), initproj.DefaultDeps(), stdout, stderr); err != nil {
		t.Fatalf("initproj.Init: %v", err)
	}

	modRoot := testModuleRoot(t)
	goModPath := filepath.Join(dir, ".wat", "go.mod")
	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	replaceLine := "\nreplace github.com/sviatsviatsviat/wat => " + modRoot + "\n"
	if err := os.WriteFile(goModPath, append(goModBytes, []byte(replaceLine)...), 0o644); err != nil {
		t.Fatalf("write go.mod replace: %v", err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = filepath.Join(dir, ".wat")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	return dir
}

func testModuleRoot(t *testing.T) string {
	t.Helper()

	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		t.Fatalf("go env GOMOD: %v", err)
	}
	mod := strings.TrimSpace(string(out))
	if mod == "" || mod == "/dev/null" {
		t.Fatal("not inside a Go module")
	}
	return filepath.Dir(mod)
}
