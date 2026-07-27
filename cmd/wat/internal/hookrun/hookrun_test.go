package hookrun_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookrun"
)

func TestHintArgs(t *testing.T) {
	t.Parallel()
	if got := hookrun.HintArgs("", ""); len(got) != 0 {
		t.Fatalf("empty = %#v", got)
	}
	if got := hookrun.HintArgs("claude", ""); !reflect.DeepEqual(got, []string{"--agent", "claude"}) {
		t.Fatalf("agent only = %#v", got)
	}
	if got := hookrun.HintArgs("", "PreToolUse"); !reflect.DeepEqual(got, []string{"--event", "PreToolUse"}) {
		t.Fatalf("event only = %#v", got)
	}
	if got := hookrun.HintArgs("cursor", "sessionStart"); !reflect.DeepEqual(got, []string{"--agent", "cursor", "--event", "sessionStart"}) {
		t.Fatalf("both = %#v", got)
	}
}

func TestRun_forwardsHintArgs(t *testing.T) {
	dir := t.TempDir()
	watDir := filepath.Join(dir, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var gotArgs []string
	deps := hookrun.DefaultDeps()
	deps.Getenv = func(string) string { return "" }
	deps.Getwd = func() (string, error) { return dir, nil }
	deps.Command = func(name string, args ...string) *exec.Cmd {
		if name == "go" {
			switch {
			case len(args) >= 1 && args[0] == "env":
				return exec.Command("echo", "go1.26.0")
			case len(args) >= 1 && args[0] == "list":
				return exec.Command("echo", "wat-hooks")
			case len(args) > 0 && args[0] == "build":
				for i := 0; i+1 < len(args); i++ {
					if args[i] == "-o" {
						if err := os.WriteFile(args[i+1], []byte("#!/bin/true\n"), 0o755); err != nil {
							t.Fatalf("stage fake binary: %v", err)
						}
						break
					}
				}
				return exec.Command("true")
			}
		}
		gotArgs = append([]string{name}, args...)
		return exec.Command("true")
	}
	deps.RunCmd = func(*exec.Cmd) error { return nil }

	code := hookrun.Run(hookrun.Config{Agent: "claude", Event: "PreToolUse"}, "vtest", deps, os.Stderr)
	if code != hookrun.ExitOK {
		t.Fatalf("exit = %d", code)
	}
	if len(gotArgs) < 5 || gotArgs[1] != "--agent" || gotArgs[2] != "claude" || gotArgs[3] != "--event" || gotArgs[4] != "PreToolUse" {
		t.Fatalf("hooks argv = %#v", gotArgs)
	}
}
