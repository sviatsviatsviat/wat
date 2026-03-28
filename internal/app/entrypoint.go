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
// ([cli.ExitSuccess], [cli.ExitGeneral], or [cli.ExitBadInput]). After the host is parsed and a [core.HookAdapterFactory] is ready,
// stdin is read and a [core.HookAdapter] is built before the [core.HookHandlerProvider] is constructed from the remaining arguments.
func Execute(programArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	console := cli.NewConsole(stderr, stdout)

	if len(programArgs) == 0 {
		cli.PrintRootHelp(console)
		return cli.ExitBadInput
	}

	hookAdapterFactory, argsAfterHost, hostFactoryReady := prepareHookAdapterFactory(console, programArgs)
	if !hostFactoryReady {
		return cli.ExitBadInput
	}

	hookAdapter, hookAdapterReady := prepareHookAdapter(console, stdin, hookAdapterFactory)
	if !hookAdapterReady {
		return cli.ExitGeneral
	}

	hookHandlerProvider, hookHandlerProviderReady := prepareHookHandlerProvider(console, argsAfterHost)
	if !hookHandlerProviderReady {
		return cli.ExitBadInput
	}

	hookHandler, hookHandlerErr := hookHandlerProvider.HookHandlerFor(hookAdapter)
	if hookHandlerErr != nil {
		_ = console.WriteError(hookHandlerErr.Error())
		return cli.ExitGeneral
	}
	result := hookHandler.Handle()
	return result.Code
}

// prepareHookAdapterFactory parses the hook host from programArgs, builds the host-specific [core.HookAdapterFactory],
// and returns the remaining arguments with the host token removed for subcommand parsing. On failure it writes to console:
// host parse errors also print root help. If ok is false, the caller should return [cli.ExitBadInput].
func prepareHookAdapterFactory(console cli.Console, programArgs []string) (hookAdapterFactory core.HookAdapterFactory, argsAfterHost []string, ok bool) {
	hookHostName, argsAfterHost, parseHostErr := parseHost(programArgs)
	if parseHostErr != nil {
		_ = console.WriteError(parseHostErr.Error())
		cli.PrintRootHelp(console)
		return nil, nil, false
	}
	hostHookAdapterFactory, newHookAdapterFactoryErr := newHookAdapterFactory(hookHostName)
	if newHookAdapterFactoryErr != nil {
		_ = console.WriteError(newHookAdapterFactoryErr.Error())
		return nil, nil, false
	}
	return hostHookAdapterFactory, argsAfterHost, true
}

// prepareHookHandlerProvider parses the wat subcommand from argsAfterHost and builds [core.HookHandlerProvider]. On failure it writes to console:
// subcommand parse errors also print root help; invalid --file-pattern regexp is written explicitly. If ok is false,
// the caller should return [cli.ExitBadInput].
func prepareHookHandlerProvider(
	console cli.Console,
	argsAfterHost []string,
) (hookHandlerProvider core.HookHandlerProvider, ok bool) {
	watSubcommand, subcommandArgs, parseSubcommandErr := parseSubcommand(argsAfterHost)
	if parseSubcommandErr != nil {
		_ = console.WriteError(parseSubcommandErr.Error())
		cli.PrintRootHelp(console)
		return nil, false
	}
	providerImpl, newHookHandlerProviderErr := newHookHandlerProvider(watSubcommand, console, subcommandArgs)
	if newHookHandlerProviderErr != nil {
		if strings.HasPrefix(newHookHandlerProviderErr.Error(), "invalid --file-pattern regexp") {
			_ = console.WriteError(newHookHandlerProviderErr.Error())
		}
		return nil, false
	}
	return providerImpl, true
}

// prepareHookAdapter reads hook event JSON from stdin and builds [core.HookAdapter] via the factory. On failure it writes
// to stderr; if ok is false, the caller should return [cli.ExitGeneral].
func prepareHookAdapter(
	console cli.Console,
	hookStdin io.Reader,
	hookAdapterFactory core.HookAdapterFactory,
) (hookAdapter core.HookAdapter, ok bool) {
	hookEventJSON, readHookStdinJSONErr := cli.ReadHookStdinJSON(hookStdin)
	if readHookStdinJSONErr != nil {
		_ = console.WriteErrorf("failed to parse stdin event JSON: %v\n", readHookStdinJSONErr)
		return nil, false
	}
	builtHookAdapter, hookAdapterFromJSONErr := hookAdapterFactory.HookAdapterFromJSON(hookEventJSON, console)
	if hookAdapterFromJSONErr != nil {
		_ = console.WriteError(hookAdapterFromJSONErr.Error())
		return nil, false
	}
	return builtHookAdapter, true
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

// parseSubcommand reads the wat subcommand from the first element of argsAfterHost; subcommandArgs are passed to the subcommand implementation (e.g. exec flags and templated command arguments).
func parseSubcommand(argsAfterHost []string) (watSubcommand string, subcommandArgs []string, err error) {
	return splitFirstArg(argsAfterHost,
		"expected wat <host> <command> … (missing command after host)",
		"command cannot be empty",
	)
}
