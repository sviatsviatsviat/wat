package e2e_test

import (
	"bytes"
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
		wantOnStdout bool
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
			name:       "help_lists_version",
			args:       []string{"help"},
			wantOutput: "version",
		},
		{
			name:       "version_help",
			args:       []string{"version", "-h"},
			wantOutput: "module version",
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
			wantOnStdout: true,
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
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected non-zero exit")
				}
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					t.Fatalf("expected exit error, got %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
				}
				if code := exitErr.ExitCode(); code != tt.wantCode {
					t.Fatalf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s", code, tt.wantCode, stdout.String(), stderr.String())
				}
				got := stderr.String()
				if tt.wantOnStdout {
					got = stdout.String()
				}
				if !strings.Contains(got, tt.wantOutput) {
					stream := "stderr"
					if tt.wantOnStdout {
						stream = "stdout"
					}
					t.Fatalf("%s missing %q:\n%s", stream, tt.wantOutput, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.wantOutput) {
				t.Fatalf("stdout missing %q:\n%s", tt.wantOutput, stdout.String())
			}
		})
	}
}

func TestWat_version(t *testing.T) {
	binary := buildWat(t)

	for _, args := range [][]string{{"version"}, {"--version"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			cmd := exec.Command(binary, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				t.Fatalf("unexpected error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
			}
			got := strings.TrimSpace(stdout.String())
			if got != "v0.0.0-e2e-000000000000" {
				t.Fatalf("stdout = %q, want %q", got, "v0.0.0-e2e-000000000000")
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}
