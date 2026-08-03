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

// testdataDirName is the project-local fixture directory. buildcache excludes the
// same top-level name from the hook binary cache key.
const testdataDirName = "testdata"

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

// Init scaffolds .wat/go.mod, .wat/hooks.go, .wat/.gitignore, and starter
// fixtures under root.
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

	// Write-if-missing before the hooks.go force guard so older projects pick up
	// .gitignore on a non-force re-run without overwriting hooks.go.
	if err := writeFileIfMissing(filepath.Join(watDir, ".gitignore"), []byte(Gitignore), deps); err != nil {
		return err
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

	if err := writeStarterTestdata(watDir, deps); err != nil {
		return err
	}

	if err := goModTidy(watDir, deps, out, errOut); err != nil {
		return err
	}

	_, _ = fmt.Fprintln(out, "Initialized .wat/ hook project.")
	_, _ = fmt.Fprintln(out, "")
	_, _ = fmt.Fprintln(out, "Next steps:")
	_, _ = fmt.Fprintln(out, "  - Edit .wat/hooks.go")
	_, _ = fmt.Fprintln(out, "  - Run wat test --agent cursor --fixture .wat/testdata/cursor/session_start.json")
	_, _ = fmt.Fprintln(out, "  - Run wat install")
	_, _ = fmt.Fprintln(out, "  - Run wat doctor")
	return nil
}

func writeStarterTestdata(watDir string, deps Deps) error {
	for _, f := range starterTestdata {
		path := filepath.Join(watDir, filepath.FromSlash(f.rel))
		if err := deps.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
		}
		if err := writeFileIfMissing(path, []byte(f.body), deps); err != nil {
			return err
		}
	}
	return nil
}

func writeFileIfMissing(path string, contents []byte, deps Deps) error {
	if info, err := deps.Stat(path); err == nil {
		if info == nil || !info.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", path)
		}
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

// Gitignore is the default .wat/.gitignore scaffold body (write-if-missing).
const Gitignore = `# Hook build cache (wat run / wat install / wat test).
.cache/
`

type starterFile struct {
	rel  string
	body string
}

// starterTestdata is written under .wat/ on init when the paths are missing.
var starterTestdata = []starterFile{
	{
		rel: testdataDirName + "/cursor/session_start.json",
		body: `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "sessionStart",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "cwd": "/w",
  "is_background_agent": false,
  "composer_mode": "agent"
}
`,
	},
	{
		rel: testdataDirName + "/cursor/session_start.expect.json",
		body: `{
  "exit": 0,
  "stdout_contains": ["wat hooks are active"]
}
`,
	},
	{
		rel: testdataDirName + "/claude/session_start.json",
		body: `{
  "session_id": "s1",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/w",
  "permission_mode": "default",
  "hook_event_name": "SessionStart",
  "source": "startup"
}
`,
	},
	{
		rel: testdataDirName + "/claude/session_start.expect.json",
		body: `{
  "exit": 0,
  "stdout_contains": ["wat hooks are active"]
}
`,
	},
	{
		rel: testdataDirName + "/copilot/session_start.json",
		body: `{
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "hook_event_name": "SessionStart",
  "source": "new",
  "initial_prompt": "hi"
}
`,
	},
	{
		rel: testdataDirName + "/copilot/session_start.expect.json",
		body: `{
  "exit": 0,
  "stdout_contains": ["wat hooks are active"]
}
`,
	},
}
