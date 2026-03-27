// Package app is the CLI entrypoint: program-argument parsing, hook stdin, host handler, and exit status.
package app

import (
	"errors"
	"io"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// Execute runs wat with program arguments programArgs (excluding the binary name), hook event bytes from stdin,
// hook protocol output on stdout, and diagnostics on stderr. It returns a process exit code
// ([cli.ExitSuccess], [cli.ExitGeneral], or [cli.ExitBadInput]). After the host is parsed and a [core.HookHandlerFactory] is ready,
// stdin is read and a [core.HookHandler] is built before the [core.Command] is constructed from the remaining arguments.
func Execute(programArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	console := cli.NewConsole(stderr, stdout)

	if len(programArgs) == 0 {
		cli.PrintRootHelp(console)
		return cli.ExitBadInput
	}

	hookHandlerFactory, argsAfterHost, hostFactoryReady := prepareHookHandlerFactory(console, programArgs)
	if !hostFactoryReady {
		return cli.ExitBadInput
	}

	hookHandler, hookHandlerReady := prepareHookHandler(console, stdin, hookHandlerFactory)
	if !hookHandlerReady {
		return cli.ExitGeneral
	}

	hookCommand, hookCommandReady := prepareHookCommand(console, argsAfterHost)
	if !hookCommandReady {
		return cli.ExitBadInput
	}

	result := hookHandler.Handle(hookCommand)
	_ = console.Write(result.Output)
	return result.Code
}

// prepareHookHandlerFactory parses the hook host from programArgs, builds the host-specific [core.HookHandlerFactory],
// and returns the remaining arguments with the host token removed for subcommand parsing. On failure it writes to console:
// host parse errors also print root help. If ok is false, the caller should return [cli.ExitBadInput].
func prepareHookHandlerFactory(console cli.Console, programArgs []string) (hookHandlerFactory core.HookHandlerFactory, argsAfterHost []string, ok bool) {
	hookHostName, argsAfterHost, parseHostErr := parseHost(programArgs)
	if parseHostErr != nil {
		_ = console.WriteError(parseHostErr.Error() + "\n")
		cli.PrintRootHelp(console)
		return nil, nil, false
	}
	hostHookHandlerFactory, newHookHandlerFactoryErr := newHookHandlerFactory(hookHostName)
	if newHookHandlerFactoryErr != nil {
		_ = console.WriteError(newHookHandlerFactoryErr.Error() + "\n")
		return nil, nil, false
	}
	return hostHookHandlerFactory, argsAfterHost, true
}

// prepareHookCommand parses the wat subcommand from argsAfterHost and builds [core.Command]. On failure it writes to console:
// subcommand parse errors also print root help; invalid --file-pattern regexp is written explicitly. If ok is false,
// the caller should return [cli.ExitBadInput].
func prepareHookCommand(
	console cli.Console,
	argsAfterHost []string,
) (hookCommand core.Command, ok bool) {
	watSubcommand, subcommandArgs, parseSubcommandErr := parseSubcommand(argsAfterHost)
	if parseSubcommandErr != nil {
		_ = console.WriteError(parseSubcommandErr.Error() + "\n")
		cli.PrintRootHelp(console)
		return nil, false
	}
	hookCommandImpl, newHookCommandErr := newHookCommand(watSubcommand, console, subcommandArgs)
	if newHookCommandErr != nil {
		if strings.HasPrefix(newHookCommandErr.Error(), "invalid --file-pattern regexp") {
			_ = console.WriteError(newHookCommandErr.Error() + "\n")
		}
		return nil, false
	}
	return hookCommandImpl, true
}

// prepareHookHandler reads hook event JSON from stdin and builds [core.HookHandler] via the factory. On failure it writes
// to stderr; if ok is false, the caller should return [cli.ExitGeneral].
func prepareHookHandler(
	console cli.Console,
	hookStdin io.Reader,
	hookHandlerFactory core.HookHandlerFactory,
) (hookHandler core.HookHandler, ok bool) {
	hookEventJSON, readHookStdinJSONErr := cli.ReadHookStdinJSON(hookStdin)
	if readHookStdinJSONErr != nil {
		_ = console.WriteErrorf("failed to parse stdin event JSON: %v\n", readHookStdinJSONErr)
		return nil, false
	}
	builtHookHandler, hookHandlerFromJSONErr := hookHandlerFactory.HookHandlerFromJSON(hookEventJSON)
	if hookHandlerFromJSONErr != nil {
		_ = console.WriteError(hookHandlerFromJSONErr.Error() + "\n")
		return nil, false
	}
	return builtHookHandler, true
}

// splitFirstArg trims and returns the first element of args and the remainder, or an error when args is empty or the first token is blank after trim.
func splitFirstArg(args []string, missingFirstTokenMsg, emptyFirstTokenMsg string) (firstToken string, rest []string, err error) {
	if len(args) < 1 {
		return "", nil, errors.New(missingFirstTokenMsg)
	}
	firstToken = strings.TrimSpace(args[0])
	if firstToken == "" {
		return "", nil, errors.New(emptyFirstTokenMsg)
	}
	return firstToken, args[1:], nil
}

// parseHost reads the hook host name from the first program argument; argsAfterHost is the remaining slice for subcommand parsing.
func parseHost(programArgs []string) (hookHostName string, argsAfterHost []string, err error) {
	return splitFirstArg(programArgs,
		"expected wat <host> <command> … (missing host)",
		"host cannot be empty",
	)
}

// parseSubcommand reads the wat subcommand from the first element of argsAfterHost; subcommandArgs are passed to the subcommand implementation (e.g. run flags and templated command arguments).
func parseSubcommand(argsAfterHost []string) (watSubcommand string, subcommandArgs []string, err error) {
	return splitFirstArg(argsAfterHost,
		"expected wat <host> <command> … (missing command after host)",
		"command cannot be empty",
	)
}
