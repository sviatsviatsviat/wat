package main_test

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWat_help(t *testing.T) {
	t.Helper()

	binary := buildWat(t)

	cmd := exec.Command(binary, "help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("wat help: %v\n%s", err, out)
	}

	text := string(out)
	if !strings.Contains(text, "wat") {
		t.Fatalf("help output missing wat: %q", text)
	}
}

func TestWat_unknownCommand(t *testing.T) {
	t.Helper()

	binary := buildWat(t)

	cmd := exec.Command(binary, "nosuchcommand")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("wat nosuchcommand: expected non-zero exit")
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("wat nosuchcommand: expected exit error, got %v\n%s", err, out)
	}
	if code := exitErr.ExitCode(); code != 1 {
		t.Fatalf("exit code = %d, want 1\n%s", code, out)
	}

	text := string(out)
	if !strings.Contains(text, "wat") {
		t.Fatalf("usage output missing wat: %q", text)
	}
}

func buildWat(t *testing.T) string {
	t.Helper()

	root := moduleRoot(t)
	dir := t.TempDir()
	name := "wat"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(dir, name)

	cmd := exec.Command("go", "build", "-o", binary, "./cmd/wat")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return binary
}

func moduleRoot(t *testing.T) string {
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
