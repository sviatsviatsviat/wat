package app

import (
	"errors"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/commands"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// errHookCommandBadInput is returned when the subcommand is unknown.
var errHookCommandBadInput = errors.New("hook command: bad input")

// newHookCommand builds the [core.Command] for subcommand using subcommandArgs; unknown
// subcommands print help and return errHookCommandBadInput.
func newHookCommand(subcommand string, console cli.Console, subprocessRunner watexec.SubprocessRunner, subcommandArgs []string) (core.Command, error) {
	switch subcommand {
	case "run":
		return commands.NewRunCommand(console, subprocessRunner, subcommandArgs)
	default:
		_ = console.WriteErrorf("wat: unknown command %q\n\n", subcommand)
		cli.PrintRootHelp(console)
		return nil, errHookCommandBadInput
	}
}
