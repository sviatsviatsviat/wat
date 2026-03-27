package run

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

// newRunFlagSet registers -f/--file-pattern; filePatternFlagTarget receives the parsed value.
func newRunFlagSet(filePatternFlagTarget *string) *flag.FlagSet {
	runFlagSet := flag.NewFlagSet("run", flag.ContinueOnError)
	runFlagSet.SetOutput(io.Discard)
	runFlagSet.StringVar(filePatternFlagTarget, "file-pattern", defaultFilePatternFlagValue, "optional regexp for FILE_PATH filter on afterFileEdit; default * means no filter")
	runFlagSet.StringVar(filePatternFlagTarget, "f", defaultFilePatternFlagValue, "shorthand for file-pattern")
	return runFlagSet
}

// parseRunArgs parses run flags from programArgs and returns the subprocess argument template and the file-pattern flag string.
// Flag or validation failures write to console and print run help where appropriate.
func parseRunArgs(console cli.Console, programArgs []string) (argsTemplate []string, filePatternFromFlags string, err error) {
	var filePatternFlagValue string
	runFlagSet := newRunFlagSet(&filePatternFlagValue)
	if parseErr := runFlagSet.Parse(programArgs); parseErr != nil {
		_ = console.WriteError(parseErr.Error() + "\n")
		cli.PrintRunHelp(console)
		return nil, "", fmt.Errorf("run: %w", parseErr)
	}
	if filePatternFlagValue == "" {
		_ = console.WriteError("file-pattern value cannot be empty\n")
		cli.PrintRunHelp(console)
		return nil, "", fmt.Errorf("file-pattern value cannot be empty")
	}
	argsTemplate = runFlagSet.Args()
	if len(argsTemplate) == 0 {
		_ = console.WriteError("missing command to run (arguments for the subprocess, e.g. go version)\n")
		cli.PrintRunHelp(console)
		return nil, "", fmt.Errorf("missing command to run")
	}
	return argsTemplate, filePatternFlagValue, nil
}

// compileRunFilePattern returns nil when filePatternFromFlags is the default (no filter); otherwise compiles it as a Go regexp.
func compileRunFilePattern(filePatternFromFlags string) (*regexp.Regexp, error) {
	if filePatternFromFlags == defaultFilePatternFlagValue {
		return nil, nil
	}
	compiledFilePathFilter, err := regexp.Compile(filePatternFromFlags)
	if err != nil {
		return nil, fmt.Errorf("invalid --file-pattern regexp: %w", err)
	}
	return compiledFilePathFilter, nil
}
