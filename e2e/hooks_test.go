package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWatTest_preToolDeny(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	tests := []struct {
		name       string
		agent      string
		fixture    string
		wantExit   int
		wantOutput []string
	}{
		{
			name:     "claude",
			agent:    "claude",
			fixture:  fixturePath(t, "claude", "pre_tool_force_push.json"),
			wantExit: 0,
			wantOutput: []string{
				"event: PreToolUse",
				"deny",
				"force pushes are not allowed",
			},
		},
		{
			name:     "copilot",
			agent:    "copilot",
			fixture:  fixturePath(t, "copilot", "pre_tool_force_push.json"),
			wantExit: 0,
			wantOutput: []string{
				`"permission_decision":"deny"`,
				"force pushes are not allowed",
			},
		},
		{
			name:     "cursor",
			agent:    "cursor",
			fixture:  fixturePath(t, "cursor", "before_shell_force_push.json"),
			wantExit: 2, // cursor.PermissionDenyExit
			wantOutput: []string{
				`"permission":"deny"`,
				"force pushes are not allowed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := runWat(t, binary, project, "test", "--agent", tt.agent, "--fixture", tt.fixture)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d\nstdout:\n%s\nstderr:\n%s", code, tt.wantExit, stdout, stderr)
			}
			out := stdout + stderr
			for _, want := range tt.wantOutput {
				if !strings.Contains(out, want) {
					t.Fatalf("output missing %q:\n%s", want, out)
				}
			}
		})
	}

	t.Run("unregistered event", func(t *testing.T) {
		cmd := exec.Command(binary, "test", "--agent", "claude", "--fixture", "-")
		cmd.Dir = project
		cmd.Env = append(os.Environ(),
			"WAT_PROJECT_DIR="+project,
			"PATH="+filepath.Dir(binary)+string(os.PathListSeparator)+os.Getenv("PATH"),
		)
		cmd.Stdin = strings.NewReader(`{"hook_event_name":"Notification","session_id":"s1"}`)
		var outBuf, errBuf strings.Builder
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		code := 0
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("run wat test: %v\n%s\n%s", err, outBuf.String(), errBuf.String())
			}
			code = exitErr.ExitCode()
		}
		if code != 3 {
			t.Fatalf("exit = %d, want 3\nstdout:\n%s\nstderr:\n%s", code, outBuf.String(), errBuf.String())
		}
		if !strings.Contains(errBuf.String(), "no claude Notification handler is registered") {
			t.Fatalf("stderr = %q", errBuf.String())
		}
	})
}
