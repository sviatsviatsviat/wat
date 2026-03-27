package run

import (
	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// NewRunCommand parses optional -f/--file-pattern from programArgs, then returns a [core.Command] that
// expands the remaining argument template with hook placeholders and runs it as a subprocess (child stderr
// via [cli.Console.ConnectErrorsFrom]).
// Empty arguments after flags writes diagnostics and run help and returns a non-nil error.
func NewRunCommand(console cli.Console, programArgs []string) (core.Command, error) {
	argsTemplate, filePatternFromFlags, err := parseRunArgs(console, programArgs)
	if err != nil {
		return nil, err
	}
	filePathFilter, err := compileRunFilePattern(filePatternFromFlags)
	if err != nil {
		return nil, err
	}
	return runCommand{
		argsTemplate:         argsTemplate,
		filePathFilterRegexp: filePathFilter,
		console:              console,
	}, nil
}
