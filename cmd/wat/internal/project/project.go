package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Layout names for a wat hook project.
const (
	// ProjectDirEnv overrides project-root discovery when set.
	ProjectDirEnv = "WAT_PROJECT_DIR"
	// DirName is the hook project directory under a workspace root.
	DirName = ".wat"
	// HooksFile is the required hook script inside DirName.
	HooksFile = "hooks.go"
	// GoModFile is the required module file inside DirName.
	GoModFile = "go.mod"
)

// Deps holds injectable filesystem dependencies for project discovery.
type Deps struct {
	Getenv func(string) string
	Getwd  func() (string, error)
	Stat   func(string) (os.FileInfo, error)
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Getenv: os.Getenv,
		Getwd:  os.Getwd,
		Stat:   os.Stat,
	}
}

// Dir returns the .wat directory path under root.
func Dir(root string) string {
	return filepath.Join(root, DirName)
}

// Resolve finds the .wat/ directory using ProjectDirEnv when set, otherwise
// walking upward from the current working directory.
func Resolve(deps Deps) (string, error) {
	root := strings.TrimSpace(deps.Getenv(ProjectDirEnv))
	if root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			return "", fmt.Errorf("resolve %s: %w", ProjectDirEnv, err)
		}
		return FromRoot(abs, deps)
	}
	cwd, err := deps.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	return FindFrom(cwd, deps)
}

// FromRoot returns the .wat/ directory under root after validating required files.
func FromRoot(root string, deps Deps) (string, error) {
	watDir := Dir(root)
	if err := MustHaveFiles(watDir, deps); err != nil {
		return "", fmt.Errorf("%s is not a wat hook project: %w", watDir, err)
	}
	return watDir, nil
}

// FindFrom walks upward from start looking for a valid .wat/ project.
func FindFrom(start string, deps Deps) (string, error) {
	dir := start
	for {
		watDir := Dir(dir)
		if err := MustHaveFiles(watDir, deps); err == nil {
			return watDir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no %s/ project found from %s (run \"wat init\" first)", DirName, start)
		}
		dir = parent
	}
}

// MustHaveFiles reports whether watDir contains HooksFile and GoModFile as regular files.
func MustHaveFiles(watDir string, deps Deps) error {
	for _, name := range []string{HooksFile, GoModFile} {
		path := filepath.Join(watDir, name)
		info, err := deps.Stat(path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", path)
		}
	}
	return nil
}
