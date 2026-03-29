package execcommand

import (
	"flag"
	"fmt"
	"io"
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

// defaultFilePatternFlagValue is the StringVar default for -f/--file-pattern. It is not compiled as
// a regexp: after Parse, this value yields a nil file pattern (no filter). Pass -f/--file-pattern
// with any other value to filter (Go regexp syntax).
const defaultFilePatternFlagValue = "*"

// newExecFlagSet registers -f/--file-pattern; filePatternFlagTarget receives the parsed value.
func newExecFlagSet(filePatternFlagTarget *string) *flag.FlagSet {
	execFlagSet := flag.NewFlagSet("exec", flag.ContinueOnError)
	execFlagSet.SetOutput(io.Discard)
	execFlagSet.StringVar(filePatternFlagTarget, "file-pattern", defaultFilePatternFlagValue, "optional regexp for FILE_PATH filter on afterFileEdit and afterTabFileEdit; default * means no filter")
	execFlagSet.StringVar(filePatternFlagTarget, "f", defaultFilePatternFlagValue, "shorthand for file-pattern")
	return execFlagSet
}

// parseExecArgs parses exec flags from programArgs and returns the subprocess argument template and the file-pattern flag string.
// Flag or validation failures write to console and print exec help where appropriate.
func parseExecArgs(console cli.Console, programArgs []string) (argsTemplate []string, filePatternFromFlags string, err error) {
	var filePatternFlagValue string
	execFlagSet := newExecFlagSet(&filePatternFlagValue)
	if parseErr := execFlagSet.Parse(programArgs); parseErr != nil {
		_ = console.WriteError(parseErr.Error() + "\n")
		cli.PrintExecHelp(console)
		return nil, "", fmt.Errorf("exec: %w", parseErr)
	}
	if filePatternFlagValue == "" {
		_ = console.WriteError("file-pattern value cannot be empty\n")
		cli.PrintExecHelp(console)
		return nil, "", fmt.Errorf("file-pattern value cannot be empty")
	}
	argsTemplate = execFlagSet.Args()
	if len(argsTemplate) == 0 {
		_ = console.WriteError("missing subprocess command (arguments for the child process, e.g. go version)\n")
		cli.PrintExecHelp(console)
		return nil, "", fmt.Errorf("missing subprocess command")
	}
	return argsTemplate, filePatternFlagValue, nil
}

// compileExecFilePattern returns nil when filePatternFromFlags is the default (no filter); otherwise compiles it as a Go regexp.
func compileExecFilePattern(filePatternFromFlags string) (*regexp.Regexp, error) {
	if filePatternFromFlags == defaultFilePatternFlagValue {
		return nil, nil
	}
	compiledFilePathFilter, err := regexp.Compile(filePatternFromFlags)
	if err != nil {
		return nil, fmt.Errorf("invalid --file-pattern regexp: %w", err)
	}
	return compiledFilePathFilter, nil
}
