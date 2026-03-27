package app

import (
	"errors"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/run"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// errHookCommandBadInput is returned when the subcommand is unknown.
var errHookCommandBadInput = errors.New("hook command: bad input")

// newHookCommand builds the [core.Command] for watSubcommand using subcommandArgs; unknown
// subcommands print help and return errHookCommandBadInput.
func newHookCommand(watSubcommand string, console cli.Console, subprocessRunner watexec.SubprocessRunner, subcommandArgs []string) (core.Command, error) {
	switch watSubcommand {
	case "run":
		return run.NewRunCommand(console, subprocessRunner, subcommandArgs)
	default:
		_ = console.WriteErrorf("wat: unknown command %q\n\n", watSubcommand)
		cli.PrintRootHelp(console)
		return nil, errHookCommandBadInput
	}
}
