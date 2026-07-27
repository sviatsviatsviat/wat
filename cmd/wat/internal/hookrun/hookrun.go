// Package hookrun builds and executes the cached .wat/hooks binary for wat run.
package hookrun

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

// Exit codes returned by Run.
const (
	ExitOK             = 0
	ExitBuildFailed    = 1
	ExitFailClosed     = 2
	ExitRuntimeFailure = 3
)

// Config holds options for Run.
type Config struct {
	// Agent is forwarded to the hooks binary as --agent when non-empty.
	Agent string
	// Event is forwarded to the hooks binary as --event when non-empty.
	Event string
	// FailClosed maps build failures to ExitFailClosed when true.
	FailClosed bool
}

// Deps holds injectable dependencies for Run.
type Deps struct {
	Getenv    func(string) string
	Getwd     func() (string, error)
	Stat      func(string) (os.FileInfo, error)
	ReadDir   func(string) ([]os.DirEntry, error)
	ReadFile  func(string) ([]byte, error)
	MkdirAll  func(string, os.FileMode) error
	WriteFile func(string, []byte, os.FileMode) error
	Command   func(string, ...string) *exec.Cmd
	RunCmd    func(*exec.Cmd) error
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Getenv:    os.Getenv,
		Getwd:     os.Getwd,
		Stat:      os.Stat,
		ReadDir:   os.ReadDir,
		ReadFile:  os.ReadFile,
		MkdirAll:  os.MkdirAll,
		WriteFile: os.WriteFile,
		Command:   exec.Command,
		RunCmd: func(cmd *exec.Cmd) error {
			return cmd.Run()
		},
	}
}

// Run resolves the project, ensures the hooks binary, and executes it.
func Run(cfg Config, version string, deps Deps, errOut io.Writer) int {
	watDir, err := project.Resolve(project.Deps{
		Getenv: deps.Getenv,
		Getwd:  deps.Getwd,
		Stat:   deps.Stat,
	})
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: %v\n", err)
		return ExitRuntimeFailure
	}

	bc := buildcache.Adapt(deps.Getenv, deps.Stat, deps.ReadDir, deps.ReadFile, deps.MkdirAll, deps.WriteFile, deps.Command)
	binPath, err := buildcache.Ensure(watDir, version, bc, errOut)
	if err != nil {
		if errors.Is(err, buildcache.ErrBuildFailed) {
			if cfg.FailClosed {
				return ExitFailClosed
			}
			return ExitBuildFailed
		}
		_, _ = fmt.Fprintf(errOut, "wat run: %v\n", err)
		return ExitRuntimeFailure
	}

	cmd := deps.Command(binPath, HintArgs(cfg.Agent, cfg.Event)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := deps.RunCmd(cmd); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		_, _ = fmt.Fprintf(errOut, "wat run: exec %s: %v\n", binPath, err)
		return ExitRuntimeFailure
	}
	return ExitOK
}

// HintArgs builds --agent/--event argv for the hooks binary.
func HintArgs(agent, event string) []string {
	var args []string
	if agent != "" {
		args = append(args, "--agent", agent)
	}
	if event != "" {
		args = append(args, "--event", event)
	}
	return args
}
