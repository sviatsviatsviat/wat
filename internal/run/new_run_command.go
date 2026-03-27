package run

import (
	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// NewRunCommand parses optional -f/--file-pattern from programArgs, then returns a [core.Command] that
// expands the remaining argv template with hook placeholders and runs it via subprocessRunner.
// Empty argv after flags writes diagnostics and run help and returns a non-nil error.
func NewRunCommand(console cli.Console, subprocessRunner watexec.SubprocessRunner, programArgs []string) (core.Command, error) {
	argvTemplate, filePatternFromFlags, err := parseRunArgv(console, programArgs)
	if err != nil {
		return nil, err
	}
	filePathFilter, err := compileRunFilePattern(filePatternFromFlags)
	if err != nil {
		return nil, err
	}
	return runCommand{
		argvTemplate:         argvTemplate,
		filePathFilterRegexp: filePathFilter,
		console:              console,
		subprocess:           subprocessRunner,
	}, nil
}
