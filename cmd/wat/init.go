package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	if err := os.MkdirAll(cacheDir, 0o750); err != nil {
		return fmt.Errorf("create %s: %w", cacheDir, err)
	}

	hooksPath := filepath.Join(watDir, "hooks.go")
	if _, err := os.Stat(hooksPath); err == nil && !force {
		return fmt.Errorf("%s exists; re-run with --force to overwrite", hooksPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", hooksPath, err)
	}

	goModPath := filepath.Join(watDir, "go.mod")
	goModText, err := watGoMod()
	if err != nil {
		return err
	}
	if err := writeFileIfMissing(goModPath, []byte(goModText)); err != nil {
		return err
	}

	if err := goModTidy(watDir, command); err != nil {
		return err
	}

	if err := os.WriteFile(hooksPath, []byte(watHooksGo), 0o600); err != nil {
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
	if err := os.WriteFile(path, contents, 0o600); err != nil {
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

func watGoMod() (string, error) {
	version := watModuleVersionFn()
	if version == "" {
		return "", fmt.Errorf("determine wat module version (build with -buildvcs=true or use a tagged build)")
	}
	return fmt.Sprintf("module wat-hooks\n\ngo 1.26\n\nrequire github.com/sviatsviatsviat/wat %s\n", version), nil
}

const watHooksGo = `package main

import (
	"context"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func main() {
	agnostic.UseHooks().OnPreTool(func(ctx context.Context, hook agnostic.PreToolEvent, r agnostic.PreToolResults) (agnostic.PreToolResult, error) {
		// Guard: block force pushes, escalate other git pushes to the user.
		// Fires on PreToolUse (Claude/Copilot) and on preToolUse /
		// beforeShellExecution (Cursor); hook.Tool.Shell is the extracted command.
		if hook.Tool == nil {
			return nil, nil
		}
		cmd := hook.Tool.Shell
		switch {
		case strings.Contains(cmd, "push --force"), strings.Contains(cmd, "push -f"):
			// → Claude: permissionDecision:"deny"; Copilot: permission_decision:"deny";
			//   Cursor: permission:"deny" + agent_message.
			return r.Deny("force pushes are not allowed; rebase and push normally"), nil
		case strings.HasPrefix(cmd, "git push"):
			// → "ask" where supported; Copilot cloud agent downgrades ask to deny.
			return r.Ask("agent wants to push to the remote"), nil
		}
		return nil, nil // zero = no opinion, default flow
	}).
		OnPostTool(func(ctx context.Context, hook agnostic.PostToolEvent, r agnostic.PostToolResults) (agnostic.PostToolResult, error) {
			// Command: after any file edit, tell the model which test command applies.
			// → additionalContext / additional_context on each dialect.
			if hook.Tool == nil {
				return nil, nil
			}
			if hook.Tool.Name == tools.ToolEdit || hook.Tool.Name == tools.ToolWrite {
				return r.Context("Run go test ./... to verify this change."), nil
			}
			return nil, nil
		}).
		OnStop(func(ctx context.Context, hook agnostic.StopEvent, r agnostic.StopResults) (agnostic.StopResult, error) {
			// Stop gate: refuse to finish the turn while the build is red.
			// → Claude/Copilot: decision:"block"+reason; Cursor: followup_message.
			// Loop guards differ per agent, so check both before re-blocking.
			if hook.Turn == nil || hook.Turn.StopHookActive || hook.Turn.LoopCount > 2 {
				return nil, nil // already retried; let it stop
			}
			if err := exec.CommandContext(ctx, "go", "build", "./...").Run(); err != nil {
				return r.FollowUp("go build ./... fails; fix the build before finishing"), nil
			}
			return nil, nil
		})

	run.Main() // auto-detects dialect, dispatches, merges, exits
}
`
