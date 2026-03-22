package app

import (
	"fmt"
	"strings"
)

// parseCommandAndCommonParameters returns the subcommand from args[0], the resolved host, and
// remaining argv after leading --host/-H pairs following the subcommand.
func parseCommandAndCommonParameters(args []string) (subcommand string, host string, subcommandArgs []string, err error) {
	if len(args) == 0 {
		return "", "", nil, fmt.Errorf("internal: empty args")
	}
	subcommand = args[0]
	host, subcommandArgs, err = consumeSubcommandHostFlags(args[1:])
	return subcommand, host, subcommandArgs, err
}

// consumeSubcommandHostFlags parses leading --host or -H flags in afterSubcommand
// (--host value, -H value, --host=value). Short -H does not accept -H=value.
// The default host name is "cursor"; later pairs override earlier ones.
func consumeSubcommandHostFlags(afterSubcommand []string) (host string, subcommandArgs []string, err error) {
	host = "cursor"
	subcommandArgs = afterSubcommand
	for {
		if len(subcommandArgs) == 0 {
			return host, subcommandArgs, nil
		}
		flagOrArg := subcommandArgs[0]
		if hostValue, ok := strings.CutPrefix(flagOrArg, "--host="); ok {
			if hostValue == "" {
				return "", nil, fmt.Errorf("host value cannot be empty after %s", flagOrArg)
			}
			host = hostValue
			subcommandArgs = subcommandArgs[1:]
			continue
		}
		if flagOrArg != "--host" && flagOrArg != "-H" {
			return host, subcommandArgs, nil
		}
		if len(subcommandArgs) < 2 {
			return "", nil, fmt.Errorf("missing host value after %s", flagOrArg)
		}
		hostValue := subcommandArgs[1]
		if hostValue == "" {
			return "", nil, fmt.Errorf("host value cannot be empty after %s", flagOrArg)
		}
		host = hostValue
		subcommandArgs = subcommandArgs[2:]
	}
}
