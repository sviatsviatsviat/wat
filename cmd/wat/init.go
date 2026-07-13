package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

func newInitCmd() *subcommandRunner {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	force := fs.Bool("force", false, "overwrite existing .wat/hooks.go")
	return &subcommandRunner{
		name:    "init",
		summary: "scaffold a .wat/ hook project",
		long: "Create .wat/go.mod and .wat/hooks.go in the current working directory.\n\n" +
			"Safe to re-run: existing files are preserved, except .wat/hooks.go which requires --force to overwrite.",
		fs: fs,
		run: func() int {
			args := fs.Args()
			if len(args) > 1 {
				_, _ = fmt.Fprintln(stderr, "wat init: expected at most one optional path argument")
				return exitUsage
			}
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			if err := initProject(root, *force, execCommand); err != nil {
				_, _ = fmt.Fprintf(stderr, "wat init: %v\n", err)
				return exitRuntimeFailure
			}
			return exitOK
		},
	}
}

var execCommand = exec.Command

func initProject(root string, force bool, command func(string, ...string) *exec.Cmd) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	watDir := filepath.Join(absRoot, ".wat")
	cacheDir := filepath.Join(watDir, ".cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create %s: %w", cacheDir, err)
	}

	hooksPath := filepath.Join(watDir, "hooks.go")
	if _, err := os.Stat(hooksPath); err == nil && !force {
		return fmt.Errorf("%s exists; re-run with --force to overwrite", hooksPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", hooksPath, err)
	}

	goModPath := filepath.Join(watDir, "go.mod")
	if err := writeFileIfMissing(goModPath, []byte(watGoMod())); err != nil {
		return err
	}

	if err := goModTidy(watDir, command); err != nil {
		return err
	}

	if err := os.WriteFile(hooksPath, []byte(watHooksGo), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", hooksPath, err)
	}

	_, _ = fmt.Fprintln(stdout, "Initialized .wat/ hook project.")
	_, _ = fmt.Fprintln(stdout, "")
	_, _ = fmt.Fprintln(stdout, "Next steps:")
	_, _ = fmt.Fprintln(stdout, "  - Edit .wat/hooks.go")
	_, _ = fmt.Fprintln(stdout, "  - Run wat install")
	_, _ = fmt.Fprintln(stdout, "  - Run wat doctor")
	return nil
}

func writeFileIfMissing(path string, contents []byte) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func goModTidy(dir string, command func(string, ...string) *exec.Cmd) error {
	cmd := command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy in %s: %w", dir, err)
	}
	return nil
}

func watGoMod() string {
	version := watModuleVersion()
	if version == "" {
		// Last-resort fallback: allow Go tooling to resolve a version implicitly.
		// This is uncommon when build metadata is present (Go builds default to -buildvcs=true).
		return "module wat-hooks\n\ngo 1.26\n"
	}
	return fmt.Sprintf("module wat-hooks\n\ngo 1.26\n\nrequire github.com/sviatsviatsviat/wat %s\n", version)
}

func watModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if v := strings.TrimSpace(info.Main.Version); v != "" && v != "(devel)" {
		return v
	}

	var revision string
	var t time.Time
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			parsed, err := time.Parse(time.RFC3339, s.Value)
			if err == nil {
				t = parsed.UTC()
			}
		}
	}
	if revision == "" || t.IsZero() {
		return ""
	}
	short := revision
	if len(short) > 12 {
		short = short[:12]
	}
	return fmt.Sprintf("v0.0.0-%s-%s", t.Format("20060102150405"), short)
}

const watHooksGo = `package main

import (
	"context"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

func main() {
	mux := agenthooks.NewMux()

	// Guard: block force pushes, escalate other git pushes to the user.
	// Fires on PreToolUse (Claude/Copilot) and on preToolUse /
	// beforeShellExecution (Cursor); ev.Tool.Shell is the extracted command.
	mux.On(agenthooks.KindPreTool, func(ctx context.Context, ev *agenthooks.Event) (agenthooks.Result, error) {
		cmd := ev.Tool.Shell
		switch {
		case strings.Contains(cmd, "push --force"), strings.Contains(cmd, "push -f"):
			// → Claude: permissionDecision:"deny"; Copilot: permissionDecision:"deny";
			//   Cursor: permission:"deny" + agent_message.
			return agenthooks.Deny("force pushes are not allowed; rebase and push normally"), nil
		case strings.HasPrefix(cmd, "git push"):
			// → "ask" where supported; Copilot cloud agent downgrades ask to deny.
			return agenthooks.Ask("agent wants to push to the remote"), nil
		}
		return agenthooks.Result{}, nil // zero Result = no opinion, default flow
	})

	// Command: after any file edit, tell the model which test command applies.
	// → additionalContext / additional_context on each dialect.
	mux.On(agenthooks.KindPostTool, func(ctx context.Context, ev *agenthooks.Event) (agenthooks.Result, error) {
		if ev.Tool.Name == agenthooks.ToolEdit || ev.Tool.Name == agenthooks.ToolWrite {
			return agenthooks.Context("Run go test ./... to verify this change."), nil
		}
		return agenthooks.Result{}, nil
	})

	// Stop gate: refuse to finish the turn while the build is red.
	// → Claude/Copilot: decision:"block"+reason; Cursor: followup_message.
	// Loop guards differ per agent, so check both before re-blocking.
	mux.On(agenthooks.KindStop, func(ctx context.Context, ev *agenthooks.Event) (agenthooks.Result, error) {
		if ev.Turn.StopHookActive || ev.Turn.LoopCount > 2 {
			return agenthooks.Result{}, nil // already retried; let it stop
		}
		if err := exec.CommandContext(ctx, "go", "build", "./...").Run(); err != nil {
			return agenthooks.Result{FollowUp: "go build ./... fails; fix the build before finishing"}, nil
		}
		return agenthooks.Result{}, nil
	})

	mux.Main() // stdin → detect agent → decode → dispatch → encode → exit
}
`
