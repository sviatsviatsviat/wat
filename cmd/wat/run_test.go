package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookrun"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { return "x" }
func (fakeFileInfo) Size() int64        { return 1 }
func (fakeFileInfo) Mode() os.FileMode  { return 0o644 }
func (fakeFileInfo) ModTime() time.Time { return time.Now() }
func (fakeFileInfo) IsDir() bool        { return false }
func (fakeFileInfo) Sys() any           { return nil }

func TestRunHook_walkUpFindsWatDir(t *testing.T) {
	dir := t.TempDir()
	proj := filepath.Join(dir, "project")
	watDir := filepath.Join(proj, ".wat")
	subdir := filepath.Join(proj, "a", "b")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package main\nfunc main(){}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	deps := hookrun.DefaultDeps()
	deps.Getenv = func(string) string { return "" }
	deps.Getwd = func() (string, error) { return subdir, nil }
	deps.Command = func(string, ...string) *exec.Cmd { return exec.Command("go", "version") }
	deps.RunCmd = func(*exec.Cmd) error { return nil }

	if got := hookrun.Run(hookrun.Config{}, "vtest", deps, stderr); got != hookrun.ExitOK {
		t.Fatalf("exit = %d, want %d", got, hookrun.ExitOK)
	}
}

func TestRunHook_buildFailureExitCode(t *testing.T) {
	var errBuf bytes.Buffer
	prevStderr := stderr
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	deps := hookrun.DefaultDeps()
	deps.Getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/tmp"
		}
		return ""
	}
	deps.Getwd = func() (string, error) { return "/tmp", nil }
	deps.Stat = func(path string) (os.FileInfo, error) {
		if strings.HasSuffix(path, filepath.Join(".wat", "hooks.go")) || strings.HasSuffix(path, filepath.Join(".wat", "go.mod")) {
			return fakeFileInfo{}, nil
		}
		return nil, os.ErrNotExist
	}
	deps.ReadDir = func(string) ([]os.DirEntry, error) { return nil, nil }
	deps.ReadFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.MkdirAll = func(string, os.FileMode) error { return nil }
	deps.WriteFile = func(string, []byte, os.FileMode) error { return nil }
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" && len(args) >= 2 && args[0] == "env" {
			return exec.Command("echo", "go1.26.0")
		}
		// buildHookBinary uses CombinedOutput; return a command guaranteed to fail.
		return exec.Command("go", "not-a-command")
	}

	code := hookrun.Run(hookrun.Config{FailClosed: false}, "vtest", deps, stderr)
	if code != hookrun.ExitBuildFailed {
		t.Fatalf("exit = %d, want %d", code, hookrun.ExitBuildFailed)
	}

	errBuf.Reset()
	code = hookrun.Run(hookrun.Config{FailClosed: true}, "vtest", deps, stderr)
	if code != hookrun.ExitFailClosed {
		t.Fatalf("exit = %d, want %d", code, hookrun.ExitFailClosed)
	}
}

func BenchmarkRunHook_cacheHitOverhead(b *testing.B) {
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	b.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	deps := hookrun.DefaultDeps()
	deps.Getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/tmp"
		}
		return ""
	}
	deps.Getwd = func() (string, error) { return "/tmp", nil }
	deps.Stat = func(string) (os.FileInfo, error) { return fakeFileInfo{}, nil }
	deps.ReadDir = func(string) ([]os.DirEntry, error) { return nil, nil }
	deps.ReadFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" {
			return exec.Command("echo", "go1.26.0")
		}
		return exec.Command("true")
	}
	deps.RunCmd = func(*exec.Cmd) error { return nil }

	cfg := hookrun.Config{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if code := hookrun.Run(cfg, "vtest", deps, stderr); code != hookrun.ExitOK {
			b.Fatalf("exit = %d", code)
		}
	}
}

func TestResolveWatDir_envOverrideErrors(t *testing.T) {
	deps := project.DefaultDeps()
	deps.Getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/does/not/exist"
		}
		return ""
	}
	deps.Stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	_, err := project.Resolve(deps)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ".wat") && !strings.Contains(err.Error(), "WAT_PROJECT_DIR") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestRunHook_execExitCodePassthrough(t *testing.T) {
	var errBuf bytes.Buffer
	prevStderr := stderr
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	deps := hookrun.DefaultDeps()
	deps.Getenv = func(string) string { return "" }
	deps.Getwd = func() (string, error) { return "/tmp", nil }
	deps.Stat = func(string) (os.FileInfo, error) { return fakeFileInfo{}, nil }
	deps.ReadDir = func(string) ([]os.DirEntry, error) { return nil, nil }
	deps.ReadFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" {
			return exec.Command("echo", "go1.26.0")
		}
		cmd := exec.Command(os.Args[0], "-test.run=^TestRunHookExecExitHelper$")
		cmd.Env = append(os.Environ(), "WAT_TEST_HELPER_EXIT=7")
		return cmd
	}
	deps.RunCmd = func(cmd *exec.Cmd) error { return cmd.Run() }

	code := hookrun.Run(hookrun.Config{}, "vtest", deps, stderr)
	if code != 7 {
		t.Fatalf("exit = %d, want 7", code)
	}

	deps.RunCmd = func(*exec.Cmd) error { return errors.New("boom") }
	code = hookrun.Run(hookrun.Config{}, "vtest", deps, stderr)
	if code != hookrun.ExitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, hookrun.ExitRuntimeFailure)
	}
}

func TestRunHookExecExitHelper(t *testing.T) {
	code := strings.TrimSpace(os.Getenv("WAT_TEST_HELPER_EXIT"))
	if code == "" {
		t.Skip("not running as helper process")
	}
	if code != "7" {
		t.Fatalf("WAT_TEST_HELPER_EXIT = %q, want %q", code, "7")
	}
	os.Exit(7)
}
