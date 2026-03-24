// Package app is the CLI entrypoint: argv parsing, hook stdin, host handler, and exit status.
package app

import (
	"io"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// Execute runs wat with program arguments args (excluding the binary name), hook event bytes from stdin,
// hook protocol output on stdout, and diagnostics on stderr. It returns a process exit code
// ([cli.ExitSuccess], [cli.ExitGeneral], or [cli.ExitBadInput]).
func Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	console := cli.NewConsole(stderr, stdout)
	subprocessRunner := watexec.NewRunner(stderr, console)

	if len(args) == 0 {
		cli.PrintRootHelp(console)
		return cli.ExitBadInput
	}

	watExecCtx, subcommandArgs, initErr := initializeContext(args)
	if initErr != nil {
		_ = console.WriteError(initErr.Error())
		cli.PrintRootHelp(console)
		return cli.ExitBadInput
	}

	hookHandlerFactory, hookHandlerFactoryErr := newHookHandlerFactory(watExecCtx)
	if hookHandlerFactoryErr != nil {
		_ = console.WriteError(hookHandlerFactoryErr.Error())
		return cli.ExitBadInput
	}

	hookCommand, hookCommandErr := newHookCommand(watExecCtx.Subcommand(), console, subprocessRunner, subcommandArgs)
	if hookCommandErr != nil {
		return cli.ExitBadInput
	}

	hookEventJSON, hookStdinErr := cli.ReadHookStdinJSON(stdin)
	if hookStdinErr != nil {
		_ = console.WriteErrorf("failed to parse stdin event JSON: %v\n", hookStdinErr)
		return cli.ExitGeneral
	}

	hookHandler, factoryErr := hookHandlerFactory.HookHandlerFromJSON(hookEventJSON)
	if factoryErr != nil {
		_ = console.WriteError(factoryErr.Error())
		return cli.ExitGeneral
	}

	result := hookHandler.Handle(hookCommand)
	_ = console.Write(result.Output)
	return result.Code
}
