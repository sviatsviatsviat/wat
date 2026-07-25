package e2e_test

import (
	"os"
	"path/filepath"
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
			name:    "claude",
			agent:   "claude",
			fixture: fixturePath(t, "claude", "session_start.json"),
			wantOutput: []string{
				"event: SessionStart",
				"wat hooks are active",
			},
		},
		{
			name:    "copilot",
			agent:   "copilot",
			fixture: fixturePath(t, "copilot", "session_start.json"),
			wantOutput: []string{
				"wat hooks are active",
			},
		},
		{
			name:    "cursor",
			agent:   "cursor",
			fixture: fixturePath(t, "cursor", "session_start.json"),
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

// cursorSubagentModelGateHooksGo registers a Cursor subagentStart handler that denies
// a subagent spawn pinned to a model other than "auto", mirroring the gate in this
// repository's own .wat/hooks.go. This exercises wat's CLI path (registration
// manifest, dispatch, permission-deny exit code) for a Cursor permission-gating event
// in a project scaffolded by "wat init", independent of the committed .wat/ hooks.
const cursorSubagentModelGateHooksGo = `package hooks

import (
	"context"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

const autoModel = "auto"

var Hooks = []run.Hooks{
	cursor.UseHooks().SubagentStart(func(_ context.Context, hook cursor.SubagentStart, r cursor.SubagentStartResults) (cursor.PermissionOutput, error) {
		model := strings.TrimSpace(hook.SubagentModel)
		if model == "" || strings.EqualFold(model, autoModel) {
			return nil, nil
		}
		return r.Deny("model " + model + " is not pre-approved").
			WithUserMessage("wat blocked this subagent: re-run with the auto model."), nil
	}),
}
`

func TestWatTest_cursorSubagentStartModelGate(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	hooksPath := filepath.Join(project, ".wat", "hooks.go")
	if err := os.WriteFile(hooksPath, []byte(cursorSubagentModelGateHooksGo), 0o600); err != nil {
		t.Fatalf("write %s: %v", hooksPath, err)
	}

	tests := []struct {
		name       string
		fixture    string
		wantExit   int
		wantOutput []string
	}{
		{
			name:       "non-auto model denied",
			fixture:    fixturePath(t, "cursor", "subagent_start.json"),
			wantExit:   2,
			wantOutput: []string{`"permission":"deny"`, "re-run with the auto model"},
		},
		{
			name:       "auto model allowed",
			fixture:    fixturePath(t, "cursor", "subagent_start_auto.json"),
			wantExit:   0,
			wantOutput: []string{"stdout: (empty)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := runWat(t, binary, project, "test", "--agent", "cursor", "--fixture", tt.fixture)
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
}
