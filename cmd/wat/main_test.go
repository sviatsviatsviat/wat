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
			wantOutput: "The commands are:",
		},
		{
			name:       "help_lists_commands",
			args:       []string{"help"},
			wantOutput: "doctor",
		},
		{
			name:       "no_args",
			args:       nil,
			wantErr:    true,
			wantCode:   1,
			wantOutput: "Usage:",
		},
		{
			name:       "unknown_command",
			args:       []string{"nosuchcommand"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "unknown command",
		},
		{
			name:       "subcommand_help",
			args:       []string{"run", "--help"},
			wantOutput: "-agent",
		},
		{
			name:       "help_subcommand",
			args:       []string{"help", "run"},
			wantOutput: "-fail-closed",
		},
		{
			name:       "run_without_project",
			args:       []string{"run"},
			wantErr:    true,
			wantCode:   3,
			wantOutput: "no .wat/ project found",
		},
		{
			name:       "install_without_project",
			args:       []string{"install"},
			wantErr:    true,
			wantCode:   3,
			wantOutput: "no .wat/ project found",
		},
		{
			name:       "port_missing_flags",
			args:       []string{"port"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "--from is required",
		},
		{
			name:       "port_help",
			args:       []string{"port", "--help"},
			wantOutput: "-from",
		},
		{
			name:       "test_missing_fixture",
			args:       []string{"test"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "--fixture is required",
		},
		{
			name:       "test_help",
			args:       []string{"test", "--help"},
			wantOutput: "-verbose",
		},
		{
			name:       "test_help_fixture_flag",
			args:       []string{"test", "--help"},
			wantOutput: "-fixture",
		},
		{
			name:       "stub_doctor",
			args:       []string{"doctor"},
			wantErr:    true,
			wantCode:   2,
			wantOutput: "not implemented",
		},
		{
			name:       "invalid_agent",
			args:       []string{"run", "--agent", "nosuch"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "unknown agent dialect",
		},
		{
			name:       "valid_agent_without_project",
			args:       []string{"run", "--agent", "claude"},
			wantErr:    true,
			wantCode:   3,
			wantOutput: "no .wat/ project found",
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
