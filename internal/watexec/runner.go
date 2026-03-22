// Package watexec runs child processes for hook commands.
package watexec

import (
	"errors"
	"io"
	"os/exec"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

// SubprocessRunner runs a subprocess from a resolved argv slice.
type SubprocessRunner interface {
	Run(args []string) int
}

// runner implements [SubprocessRunner] with a fixed child stderr writer and console for errors.
type runner struct {
	stderr  io.Writer
	console cli.Console
}

// NewRunner returns a [SubprocessRunner]; stderr receives the child's stderr (often the same stream as console diagnostics).
func NewRunner(stderr io.Writer, console cli.Console) SubprocessRunner {
	return &runner{stderr: stderr, console: console}
}

// Run executes args[0] with args[1:] when it resolves on PATH; otherwise runs the same argv via
// the system shell (e.g. Windows builtins like echo). Child stdout is discarded; stderr is r.stderr.
// It returns the child exit code, or [cli.ExitBadInput] / [cli.ExitGeneral] on failure to start or run.
func (r *runner) Run(args []string) int {
	if len(args) == 0 || args[0] == "" {
		_ = r.console.WriteError("no command to execute after templating")
		return cli.ExitBadInput
	}

	childCmd, err := commandForArgv(args)
	if err != nil {
		_ = r.console.WriteErrorf("failed to resolve command: %v\n", err)
		return cli.ExitGeneral
	}
	childCmd.Stdout = io.Discard
	childCmd.Stderr = r.stderr

	if err = childCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		_ = r.console.WriteErrorf("failed to execute command: %v\n", err)
		return cli.ExitGeneral
	}

	return cli.ExitSuccess
}
