package initproj

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

// Deps holds injectable dependencies for Init.
type Deps struct {
	Command   func(string, ...string) *exec.Cmd
	Stat      func(string) (os.FileInfo, error)
	MkdirAll  func(string, os.FileMode) error
	WriteFile func(string, []byte, os.FileMode) error
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Command:   exec.Command,
		Stat:      os.Stat,
		MkdirAll:  os.MkdirAll,
		WriteFile: os.WriteFile,
	}
}

// Init scaffolds .wat/go.mod and .wat/hooks.go under root.
func Init(root string, force bool, version string, deps Deps, out, errOut io.Writer) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	watDir := project.Dir(absRoot)
	cacheDir := filepath.Join(watDir, ".cache")
	if err := deps.MkdirAll(cacheDir, 0o750); err != nil {
		return fmt.Errorf("create %s: %w", cacheDir, err)
	}

	hooksPath := filepath.Join(watDir, project.HooksFile)
	if _, err := deps.Stat(hooksPath); err == nil && !force {
		return fmt.Errorf("%s exists; re-run with --force to overwrite", hooksPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", hooksPath, err)
	}

	goModPath := filepath.Join(watDir, project.GoModFile)
	goModText, err := GoMod(version)
	if err != nil {
		return err
	}
	if err := writeFileIfMissing(goModPath, []byte(goModText), deps); err != nil {
		return err
	}

	if err := deps.WriteFile(hooksPath, []byte(HooksGo), 0o600); err != nil {
		return fmt.Errorf("write %s: %w", hooksPath, err)
	}

	if err := goModTidy(watDir, deps, out, errOut); err != nil {
		return err
	}

	_, _ = fmt.Fprintln(out, "Initialized .wat/ hook project.")
	_, _ = fmt.Fprintln(out, "")
	_, _ = fmt.Fprintln(out, "Next steps:")
	_, _ = fmt.Fprintln(out, "  - Edit .wat/hooks.go")
	_, _ = fmt.Fprintln(out, "  - Run wat install")
	_, _ = fmt.Fprintln(out, "  - Run wat doctor")
	return nil
}

func writeFileIfMissing(path string, contents []byte, deps Deps) error {
	if _, err := deps.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if err := deps.WriteFile(path, contents, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func goModTidy(dir string, deps Deps, out, errOut io.Writer) error {
	cmd := deps.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = errOut
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy in %s: %w", dir, err)
	}
	return nil
}

// GoMod returns the scaffolded go.mod body for the given wat module version.
func GoMod(version string) (string, error) {
	if version == "" {
		return "", fmt.Errorf("determine wat module version (build with -buildvcs=true or use a tagged build)")
	}
	return fmt.Sprintf("module wat-hooks\n\ngo 1.26\n\nrequire github.com/sviatsviatsviat/wat %s\n", version), nil
}

// HooksGo is the default .wat/hooks.go scaffold template.
const HooksGo = `package hooks

import (
	"context"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hooks contains this project's hook registrations.
var Hooks = []run.Hooks{
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
			}),
}
`
