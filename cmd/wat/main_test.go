package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func packageDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(file)
}

func TestMain_NoArgsUsageExitBadInput(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	name := "wat"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(dir, name)
	buildWatCmd := exec.Command("go", "build", "-o", bin, ".")
	buildWatCmd.Dir = packageDir(t)
	if buildOutput, err := buildWatCmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, buildOutput)
	}
	watCmd := exec.Command(bin)
	var stderr strings.Builder
	watCmd.Stderr = &stderr
	err := watCmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", exitErr.ExitCode())
	}
	if !strings.Contains(stderr.String(), "wat <command>") {
		t.Fatalf("expected help output, got %q", stderr.String())
	}
}
