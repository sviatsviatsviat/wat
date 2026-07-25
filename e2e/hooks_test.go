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
		name  string
		agent string
	}{
		{name: "claude", agent: "claude"},
		{name: "copilot", agent: "copilot"},
		{name: "cursor", agent: "cursor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := filepath.Join(project, ".wat", "testdata", tt.agent, "session_start.json")
			stdout, stderr, code := runWat(t, binary, project, "test", "--agent", tt.agent, "--fixture", fixture)
			if code != 0 {
				t.Fatalf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
			}
			out := stdout + stderr
			for _, want := range []string{"status: pass", "wat hooks are active"} {
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

const cursorSubagentDenyExpect = `{
  "exit": 2,
  "decision": "deny",
  "stdout_contains": ["re-run with the auto model"]
}
`

const cursorSubagentAllowExpect = `{
  "exit": 0,
  "stdout_empty": true
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
		name    string
		fixture string
		expect  string
	}{
		{
			name:    "non-auto model denied",
			fixture: fixturePath(t, "cursor", "subagent_start.json"),
			expect:  cursorSubagentDenyExpect,
		},
		{
			name:    "auto model allowed",
			fixture: fixturePath(t, "cursor", "subagent_start_auto.json"),
			expect:  cursorSubagentAllowExpect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectPath := filepath.Join(t.TempDir(), "case.expect.json")
			if err := os.WriteFile(expectPath, []byte(tt.expect), 0o600); err != nil {
				t.Fatal(err)
			}
			stdout, stderr, code := runWat(t, binary, project, "test", "--agent", "cursor", "--fixture", tt.fixture, "--expect", expectPath)
			if code != 0 {
				t.Fatalf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
			}
			if !strings.Contains(stdout, "status: pass") {
				t.Fatalf("output missing expect pass:\n%s\n%s", stdout, stderr)
			}
		})
	}
}

func TestWatTest_expectMismatch(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	fixture := filepath.Join(project, ".wat", "testdata", "cursor", "session_start.json")
	expectPath := filepath.Join(t.TempDir(), "wrong.expect.json")
	if err := os.WriteFile(expectPath, []byte(`{"exit":2,"decision":"deny"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := runWat(t, binary, project, "test", "--agent", "cursor", "--fixture", fixture, "--expect", expectPath)
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "expect failed") {
		t.Fatalf("stderr = %q", stderr)
	}
	if !strings.Contains(stdout, "status: fail") {
		t.Fatalf("stdout missing fail status:\n%s", stdout)
	}
}
