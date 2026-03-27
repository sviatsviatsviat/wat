package run

import (
	"errors"
	"io"
	"os/exec"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

// runSubprocess executes args[0] with args[1:] when it resolves on PATH; otherwise runs the same
// arguments via the system shell (e.g. Windows builtins like echo). Child stdout is discarded; stderr
// is wired through console.ConnectErrorsFrom. It returns the child exit code, or [cli.ExitBadInput]
// / [cli.ExitGeneral] on failure to start or run.
func runSubprocess(console cli.Console, args []string) int {
	if len(args) == 0 || args[0] == "" {
		_ = console.WriteError("no command to execute after templating")
		return cli.ExitBadInput
	}

	childCmd, err := commandForArgs(args)
	if err != nil {
		_ = console.WriteErrorf("failed to resolve command: %v\n", err)
		return cli.ExitGeneral
	}
	childCmd.Stdout = io.Discard
	console.ConnectErrorsFrom(childCmd)

	if err = childCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		_ = console.WriteErrorf("failed to execute command: %v\n", err)
		return cli.ExitGeneral
	}

	return cli.ExitSuccess
}
