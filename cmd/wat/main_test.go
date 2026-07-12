package main_test

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWat_usage(t *testing.T) {
	binary := buildWat(t)

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantCode   int
		wantOutput string
	}{
		{
			name:       "help",
			args:       []string{"help"},
			wantOutput: "Run with -h or help for this message.",
		},
		{
			name:       "unknown_command",
			args:       []string{"nosuchcommand"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "Commands will be added in later releases",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			out, err := cmd.CombinedOutput()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected non-zero exit")
				}
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("expected exit error, got %v\n%s", err, out)
				}
				if code := exitErr.ExitCode(); code != tt.wantCode {
					t.Fatalf("exit code = %d, want %d\n%s", code, tt.wantCode, out)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v\n%s", err, out)
			}

			text := string(out)
			if !strings.Contains(text, tt.wantOutput) {
				t.Fatalf("output missing %q: %q", tt.wantOutput, text)
			}
		})
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
