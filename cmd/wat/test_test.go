package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
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

	deps := defaultTestDeps()
	deps.getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return project
		}
		return ""
	}
	deps.readFixture = func(string, io.Reader) ([]byte, error) {
		return nil, nil
	}

	code := runTest(testConfig{fixture: "-"}, deps)
	if code != exitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, exitRuntimeFailure)
	}
	if !strings.Contains(errBuf.String(), "empty fixture") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestDecodeFixtureSummary_copilotRequiresEvent(t *testing.T) {
	payload, err := os.ReadFile(fixturePath(t, "copilot", "pre_tool_force_push.json"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = decodeFixtureSummary("copilot", "", payload, os.Getenv)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--event") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestWriteTestReport_eventSummary(t *testing.T) {
	ev := &agnostic.Event{
		Kind: agnostic.KindPreTool,
		Name: "PreToolUse",
		Tool: &agnostic.ToolCall{Name: agnostic.ToolBash, Shell: "git push --force"},
	}
	var buf bytes.Buffer
	writeTestReport(&buf, ev, agnostic.Claude, []byte(`{"permissionDecision":"deny","reason":"blocked"}`), nil, 0, false)

	out := buf.String()
	for _, want := range []string{"kind:  PreTool", "tool:  bash", "shell: git push --force", "decision: deny", "exit:   0"} {
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
				"kind:  PreTool",
				"deny",
				"force pushes are not allowed",
			},
		},
		{
			name:     "copilot",
			agent:    "copilot",
			event:    "preToolUse",
			fixture:  filepath.Join(fixtures, "testdata", "fixtures", "copilot", "pre_tool_force_push.json"),
			wantExit: exitOK,
			wantOutput: []string{
				`"permissionDecision":"deny"`,
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
			deps := defaultTestDeps()
			deps.getenv = func(key string) string {
				if key == "WAT_PROJECT_DIR" {
					return project
				}
				return os.Getenv(key)
			}
			deps.writeReport = &report

			code := runTest(testConfig{
				agent:   tt.agent,
				event:   tt.event,
				fixture: tt.fixture,
			}, deps)
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

	if err := initProject(dir, false, exec.Command); err != nil {
		t.Fatalf("initProject: %v", err)
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

func fixturePath(t *testing.T, parts ...string) string {
	t.Helper()
	all := append([]string{"testdata", "fixtures"}, parts...)
	return filepath.Join(testModuleRoot(t), filepath.Join(all...))
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
