package app

import (
	"errors"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/execcommand"
)

// errHookCommandBadInput is returned when the subcommand is unknown.
var errHookCommandBadInput = errors.New("hook command: bad input")

// newHookCommand builds the [core.Command] for watSubcommand using subcommandArgs; unknown
// subcommands print help and return errHookCommandBadInput.
func newHookCommand(watSubcommand string, console cli.Console, subcommandArgs []string) (core.Command, error) {
	switch watSubcommand {
	case "exec":
		return execcommand.NewExecCommand(console, subcommandArgs)
	default:
		_ = console.WriteErrorf("wat: unknown command %q\n\n", watSubcommand)
		cli.PrintRootHelp(console)
		return nil, errHookCommandBadInput
	}
}
