package e2e_test

import (
	"strings"
	"testing"
)

func TestWatTest_sessionStartContext(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	tests := []struct {
		name       string
		agent      string
		fixture    string
		wantOutput []string
	}{
		{
			name:     "claude",
			agent:    "claude",
			fixture:  fixturePath(t, "claude", "session_start.json"),
			wantOutput: []string{
				"event: SessionStart",
				"wat hooks are active",
			},
		},
		{
			name:     "copilot",
			agent:    "copilot",
			fixture:  fixturePath(t, "copilot", "session_start.json"),
			wantOutput: []string{
				"wat hooks are active",
			},
		},
		{
			name:     "cursor",
			agent:    "cursor",
			fixture:  fixturePath(t, "cursor", "session_start.json"),
			wantOutput: []string{
				"wat hooks are active",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := runWat(t, binary, project, "test", "--agent", tt.agent, "--fixture", tt.fixture)
			if code != 0 {
				t.Fatalf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
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
		stdin := strings.NewReader(`{"hook_event_name":"Notification","session_id":"s1"}`)
		_, stderr, code := runWatWithStdin(t, binary, project, stdin, "test", "--agent", "claude", "--fixture", "-")
		if code != 3 {
			t.Fatalf("exit = %d, want 3\nstderr:\n%s", code, stderr)
		}
		if !strings.Contains(stderr, "no claude Notification handler is registered") {
			t.Fatalf("stderr = %q", stderr)
		}
	})
}
