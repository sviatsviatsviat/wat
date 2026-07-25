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

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Hooks contains this project's hook registrations.
var Hooks = []run.Hooks{
	agnostic.UseHooks().OnSessionStart(func(ctx context.Context, hook agnostic.SessionStartEvent, r agnostic.SessionStartResults) (agnostic.SessionStartResult, error) {
		return r.Context("wat hooks are active"), nil
	}),
}
`
