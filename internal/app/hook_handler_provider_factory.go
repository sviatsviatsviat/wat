package app

import (
	"errors"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/execcommand"
)

// errHookHandlerProviderBadInput is returned when the subcommand is unknown.
var errHookHandlerProviderBadInput = errors.New("hook handler provider: bad input")

// newHookHandlerProvider builds the [core.HookHandlerProvider] for watSubcommand using subcommandArgs; unknown
// subcommands print help and return errHookHandlerProviderBadInput.
func newHookHandlerProvider(watSubcommand string, console cli.Console, subcommandArgs []string) (core.HookHandlerProvider, error) {
	switch watSubcommand {
	case "exec":
		return execcommand.NewExecHookHandlerProvider(console, subcommandArgs)
	default:
		_ = console.WriteErrorf("wat: unknown command %q\n\n", watSubcommand)
		cli.PrintRootHelp(console)
		return nil, errHookHandlerProviderBadInput
	}
}
