package app

import (
	"flag"
	"fmt"
	"io"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// defaultHookHost is the hook host when -H/--host is omitted.
const defaultHookHost = "cursor"

// defaultFilePatternFlagValue is the StringVar default for -f/--file-pattern. It is not compiled as
// a regexp: after Parse, this value becomes a nil FilePattern (no filter). Pass -f/--file-pattern
// with any other value to filter (Go regexp syntax).
const defaultFilePatternFlagValue = "*"

func initializeContext(args []string) (core.WatExecutionContext, []string, error) {
	if len(args) == 0 {
		return core.WatExecutionContext{}, nil, fmt.Errorf("internal: empty args")
	}
	subcommand := args[0]
	host, filePattern, subcommandArgs, err := consumeSubcommandSharedFlags(args[1:])
	if err != nil {
		return core.WatExecutionContext{}, nil, err
	}
	execCtx := core.NewWatExecutionContext(host).WithSubcommand(subcommand)
	if filePattern != nil {
		execCtx = execCtx.WithFilePattern(*filePattern)
	}
	return execCtx, subcommandArgs, nil
}

func consumeSubcommandSharedFlags(afterSubcommand []string) (host string, filePattern *string, subcommandArgs []string, err error) {
	host = defaultHookHost
	var patternStr string

	fs := flag.NewFlagSet("wat", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.StringVar(&host, "host", defaultHookHost, "hook host")
	fs.StringVar(&host, "H", defaultHookHost, "hook host (shorthand)")
	fs.StringVar(&patternStr, "file-pattern", defaultFilePatternFlagValue, "optional regexp for afterFileEdit file_path filter; default * means no filter")
	fs.StringVar(&patternStr, "f", defaultFilePatternFlagValue, "shorthand for file-pattern")

	if err := fs.Parse(afterSubcommand); err != nil {
		return "", nil, nil, err
	}

	if host == "" {
		return "", nil, nil, fmt.Errorf("host value cannot be empty after --host")
	}

	if patternStr == "" {
		return "", nil, nil, fmt.Errorf("file-pattern value cannot be empty")
	}

	if patternStr == defaultFilePatternFlagValue {
		return host, nil, fs.Args(), nil
	}
	v := patternStr
	return host, &v, fs.Args(), nil
}
