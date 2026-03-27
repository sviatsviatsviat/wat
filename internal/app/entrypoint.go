// Package app is the CLI entrypoint: argv parsing, hook stdin, host handler, and exit status.
package app

import (
	"errors"
	"io"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// Execute runs wat with program arguments programArgs (excluding the binary name), hook event bytes from stdin,
// hook protocol output on stdout, and diagnostics on stderr. It returns a process exit code
// ([cli.ExitSuccess], [cli.ExitGeneral], or [cli.ExitBadInput]). After argv yields a [core.HookHandlerFactory],
// stdin is read and a [core.HookHandler] is built before the [core.Command] is constructed from the remaining argv.
func Execute(programArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	console := cli.NewConsole(stderr, stdout)
	subprocessRunner := watexec.NewRunner(stderr, console)

	if len(programArgs) == 0 {
		cli.PrintRootHelp(console)
		return cli.ExitBadInput
	}

	hookHandlerFactory, argvAfterHost, hostFactoryReady := prepareHookHandlerFactory(console, programArgs)
	if !hostFactoryReady {
		return cli.ExitBadInput
	}

	hookHandler, hookHandlerReady := prepareHookHandler(console, stdin, hookHandlerFactory)
	if !hookHandlerReady {
		return cli.ExitGeneral
	}

	hookCommand, hookCommandReady := prepareHookCommand(console, subprocessRunner, argvAfterHost)
	if !hookCommandReady {
		return cli.ExitBadInput
	}

	result := hookHandler.Handle(hookCommand)
	_ = console.Write(result.Output)
	return result.Code
}

// prepareHookHandlerFactory parses the hook host from argv, builds the host-specific [core.HookHandlerFactory],
// and returns argv with the host token removed for subcommand parsing. On failure it writes to console:
// host parse errors also print root help. If ok is false, the caller should return [cli.ExitBadInput].
func prepareHookHandlerFactory(console cli.Console, programArgs []string) (hookHandlerFactory core.HookHandlerFactory, argvAfterHost []string, ok bool) {
	hookHostName, argvAfterHost, parseHostErr := parseHost(programArgs)
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
	return hostHookHandlerFactory, argvAfterHost, true
}

// prepareHookCommand parses the wat subcommand from argvAfterHost and builds [core.Command]. On failure it writes to console:
// subcommand parse errors also print root help; invalid --file-pattern regexp is written explicitly. If ok is false,
// the caller should return [cli.ExitBadInput].
func prepareHookCommand(
	console cli.Console,
	subprocessRunner watexec.SubprocessRunner,
	argvAfterHost []string,
) (hookCommand core.Command, ok bool) {
	watSubcommand, subcommandArgs, parseSubcommandErr := parseSubcommand(argvAfterHost)
	if parseSubcommandErr != nil {
		_ = console.WriteError(parseSubcommandErr.Error() + "\n")
		cli.PrintRootHelp(console)
		return nil, false
	}
	hookCommandImpl, newHookCommandErr := newHookCommand(watSubcommand, console, subprocessRunner, subcommandArgs)
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

// splitFirstArgvToken trims and returns the first argv element and the remainder, or an error when argv is empty or the first token is blank after trim.
func splitFirstArgvToken(argv []string, missingFirstTokenMsg, emptyFirstTokenMsg string) (firstToken string, remainingArgv []string, err error) {
	if len(argv) < 1 {
		return "", nil, errors.New(missingFirstTokenMsg)
	}
	firstToken = strings.TrimSpace(argv[0])
	if firstToken == "" {
		return "", nil, errors.New(emptyFirstTokenMsg)
	}
	return firstToken, argv[1:], nil
}

// parseHost reads the hook host name from the first argv word; argvAfterHost is the remaining slice for subcommand parsing.
func parseHost(argv []string) (hookHostName string, argvAfterHost []string, err error) {
	return splitFirstArgvToken(argv,
		"expected wat <host> <command> … (missing host)",
		"host cannot be empty",
	)
}

// parseSubcommand reads the wat subcommand from the first element of argvAfterHost; subcommandArgs are passed to the subcommand implementation (e.g. run flags and template argv).
func parseSubcommand(argvAfterHost []string) (watSubcommand string, subcommandArgs []string, err error) {
	return splitFirstArgvToken(argvAfterHost,
		"expected wat <host> <command> … (missing command after host)",
		"command cannot be empty",
	)
}
