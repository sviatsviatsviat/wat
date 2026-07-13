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
)

type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { return "x" }
func (fakeFileInfo) Size() int64        { return 1 }
func (fakeFileInfo) Mode() os.FileMode  { return 0o644 }
func (fakeFileInfo) ModTime() time.Time { return time.Now() }
func (fakeFileInfo) IsDir() bool        { return false }
func (fakeFileInfo) Sys() any           { return nil }

type fakeDirEntry struct {
	name string
}

func (e fakeDirEntry) Name() string               { return e.name }
func (e fakeDirEntry) IsDir() bool                { return false }
func (e fakeDirEntry) Type() os.FileMode          { return 0 }
func (e fakeDirEntry) Info() (os.FileInfo, error) { return fakeFileInfo{}, nil }

func TestRunHook_walkUpFindsWatDir(t *testing.T) {
	dir := t.TempDir()
	project := filepath.Join(dir, "project")
	watDir := filepath.Join(project, ".wat")
	subdir := filepath.Join(project, "a", "b")
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

	deps := defaultRunDeps()
	deps.getenv = func(string) string { return "" }
	deps.getwd = func() (string, error) { return subdir, nil }
	deps.stat = os.Stat
	deps.readFile = os.ReadFile
	deps.mkdirAll = os.MkdirAll
	deps.command = func(string, ...string) *exec.Cmd { return exec.Command("go", "version") }
	deps.runCmd = func(*exec.Cmd) error { return nil }

	if got := runHook(runConfig{}, deps); got != exitOK {
		t.Fatalf("exit = %d, want %d", got, exitOK)
	}
}

func TestCacheKey_changesOnInput(t *testing.T) {
	deps := defaultRunDeps()
	deps.readDir = func(string) ([]os.DirEntry, error) {
		return []os.DirEntry{
			fakeDirEntry{name: "a.go"},
			fakeDirEntry{name: "b.go"},
		}, nil
	}
	deps.readFile = func(path string) ([]byte, error) {
		switch filepath.Base(path) {
		case "go.mod":
			return []byte("mod"), nil
		case "go.sum":
			return []byte("sum"), nil
		case "a.go":
			return []byte("hooks-a"), nil
		case "b.go":
			return []byte("hooks-b"), nil
		default:
			return nil, os.ErrNotExist
		}
	}

	a, err := hookBuildCacheKey("/tmp/.wat", deps)
	if err != nil {
		t.Fatal(err)
	}
	deps.readFile = func(path string) ([]byte, error) {
		switch filepath.Base(path) {
		case "go.mod":
			return []byte("mod2"), nil
		case "go.sum":
			return []byte("sum"), nil
		case "a.go":
			return []byte("hooks-a"), nil
		case "b.go":
			return []byte("hooks-b"), nil
		default:
			return nil, os.ErrNotExist
		}
	}
	b, err := hookBuildCacheKey("/tmp/.wat", deps)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatalf("cache key should differ when go.mod changes")
	}
}

func TestRunHook_buildFailureExitCode(t *testing.T) {
	var errBuf bytes.Buffer
	prevStderr := stderr
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	deps := defaultRunDeps()
	deps.getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/tmp"
		}
		return ""
	}
	deps.getwd = func() (string, error) { return "/tmp", nil }
	deps.stat = func(path string) (os.FileInfo, error) {
		if strings.HasSuffix(path, filepath.Join(".wat", "hooks.go")) || strings.HasSuffix(path, filepath.Join(".wat", "go.mod")) {
			return fakeFileInfo{}, nil
		}
		return nil, os.ErrNotExist
	}
	deps.readDir = func(string) ([]os.DirEntry, error) { return nil, nil }
	deps.readFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.mkdirAll = func(string, os.FileMode) error { return nil }
	deps.command = func(name string, args ...string) *exec.Cmd {
		// buildHookBinary uses CombinedOutput; return a command guaranteed to fail.
		return exec.Command("go", "not-a-command")
	}

	code := runHook(runConfig{failClosed: false}, deps)
	if code != exitBuildFailed {
		t.Fatalf("exit = %d, want %d", code, exitBuildFailed)
	}

	errBuf.Reset()
	code = runHook(runConfig{failClosed: true}, deps)
	if code != exitFailClosed {
		t.Fatalf("exit = %d, want %d", code, exitFailClosed)
	}
}

func BenchmarkRunHook_cacheHitOverhead(b *testing.B) {
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	b.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	deps := defaultRunDeps()
	deps.getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/tmp"
		}
		return ""
	}
	deps.getwd = func() (string, error) { return "/tmp", nil }
	deps.stat = func(string) (os.FileInfo, error) { return fakeFileInfo{}, nil }
	deps.readFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.mkdirAll = func(string, os.FileMode) error { return nil }
	deps.command = func(string, ...string) *exec.Cmd { return exec.Command("go", "version") }
	deps.runCmd = func(*exec.Cmd) error { return nil }

	cfg := runConfig{agent: "claude", event: "PreToolUse"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if code := runHook(cfg, deps); code != exitOK {
			b.Fatalf("exit = %d", code)
		}
	}
}

func TestResolveWatDir_envOverrideErrors(t *testing.T) {
	deps := defaultRunDeps()
	deps.getenv = func(key string) string {
		if key == "WAT_PROJECT_DIR" {
			return "/does/not/exist"
		}
		return ""
	}
	deps.stat = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }

	_, err := resolveWatDir(deps)
	if err == nil {
		t.Fatal("expected error")
	}
	// Error may be wrapped by the implementation; assert it's actionable (mentions .wat/).
	if !strings.Contains(err.Error(), ".wat") && !strings.Contains(err.Error(), "WAT_PROJECT_DIR") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestRunHook_execExitCodePassthrough(t *testing.T) {
	var errBuf bytes.Buffer
	prevStderr := stderr
	stderr = &errBuf
	t.Cleanup(func() { stderr = prevStderr })

	deps := defaultRunDeps()
	deps.getenv = func(string) string { return "" }
	deps.getwd = func() (string, error) { return "/tmp", nil }
	deps.stat = func(string) (os.FileInfo, error) { return fakeFileInfo{}, nil }
	deps.readDir = func(string) ([]os.DirEntry, error) { return nil, nil }
	deps.readFile = func(string) ([]byte, error) { return []byte("x"), nil }
	deps.command = func(string, ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], "-test.run=^TestRunHookExecExitHelper$")
		cmd.Env = append(os.Environ(), "WAT_TEST_HELPER_EXIT=7")
		return cmd
	}
	deps.runCmd = func(cmd *exec.Cmd) error { return cmd.Run() }

	code := runHook(runConfig{}, deps)
	if code != 7 {
		t.Fatalf("exit = %d, want 7", code)
	}

	// Ensure non-ExitError failures map to exitRuntimeFailure.
	deps.runCmd = func(*exec.Cmd) error { return errors.New("boom") }
	code = runHook(runConfig{}, deps)
	if code != exitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, exitRuntimeFailure)
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
