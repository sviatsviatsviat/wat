package e2e_test

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestWat_usage(t *testing.T) {
	binary := buildWat(t)

	tests := []struct {
		name         string
		args         []string
		wantErr      bool
		wantCode     int
		wantOutput   string
		emptyWorkDir bool
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
			name:         "run_without_project",
			args:         []string{"run"},
			wantErr:      true,
			wantCode:     3,
			wantOutput:   "no .wat/ project found",
			emptyWorkDir: true,
		},
		{
			name:         "install_without_project",
			args:         []string{"install"},
			wantErr:      true,
			wantCode:     3,
			wantOutput:   "no .wat/ project found",
			emptyWorkDir: true,
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
			name:         "doctor_without_project",
			args:         []string{"doctor"},
			wantErr:      true,
			wantCode:     4,
			wantOutput:   "no .wat/ project found",
			emptyWorkDir: true,
		},
		{
			name:       "invalid_agent",
			args:       []string{"run", "--agent", "nosuch"},
			wantErr:    true,
			wantCode:   1,
			wantOutput: "unknown agent dialect",
		},
		{
			name:         "valid_agent_without_project",
			args:         []string{"run", "--agent", "claude"},
			wantErr:      true,
			wantCode:     3,
			wantOutput:   "no .wat/ project found",
			emptyWorkDir: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			if tt.emptyWorkDir {
				cmd.Dir = t.TempDir()
				cmd.Env = append(cmd.Environ(), "WAT_PROJECT_DIR=")
			}
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
