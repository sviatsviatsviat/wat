package execcommand

import (
	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// NewExecCommand parses optional -f/--file-pattern from programArgs, then returns a [core.Command] that
// expands the remaining argument template with hook placeholders and runs it as a subprocess (child stderr
// via [cli.Console.ConnectErrorsFrom]).
// Empty arguments after flags writes diagnostics and exec help and returns a non-nil error.
func NewExecCommand(console cli.Console, programArgs []string) (core.Command, error) {
	argsTemplate, filePatternFromFlags, err := parseExecArgs(console, programArgs)
	if err != nil {
		return nil, err
	}
	filePathFilter, err := compileExecFilePattern(filePatternFromFlags)
	if err != nil {
		return nil, err
	}
	return execCommand{
		argsTemplate:         argsTemplate,
		filePathFilterRegexp: filePathFilter,
		console:              console,
	}, nil
}
