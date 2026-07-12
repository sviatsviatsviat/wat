package execcommand

import (
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

type execHookHandlerProvider struct {
	argsTemplate         []string
	filePathFilterRegexp *regexp.Regexp
	console              cli.Console
}

func (p *execHookHandlerProvider) HookHandlerFor(hook core.HookAdapter) (core.HookHandler, error) {
	base := execHookHandlerBase{
		argsTemplate: p.argsTemplate,
		console:      p.console,
	}
	if h, ok := hook.AsAfterFileEdit(); ok {
		return execHookHandlerAfterFileEdit{
			execHookHandlerBase:  base,
			filePathFilterRegexp: p.filePathFilterRegexp,
			hook:                 h,
		}, nil
	}
	if h, ok := hook.AsAfterShellExecution(); ok {
		return execHookHandlerAfterShell{
			execHookHandlerBase: base,
			hook:                h,
		}, nil
	}
	return nil, core.HookAdapterNotSupportedError(hook)
}

// NewExecHookHandlerProvider parses optional -f/--file-pattern from programArgs, then returns a [core.HookHandlerProvider] whose
// handlers expand the remaining argument template with hook placeholders and run it as a subprocess (child stderr
// via [cli.Console.ConnectErrorsFrom]).
// Empty arguments after flags writes diagnostics and exec help and returns a non-nil error.
func NewExecHookHandlerProvider(console cli.Console, programArgs []string) (core.HookHandlerProvider, error) {
	argsTemplate, filePatternFromFlags, err := parseExecArgs(console, programArgs)
	if err != nil {
		return nil, err
	}
	filePathFilter, err := compileExecFilePattern(filePatternFromFlags)
	if err != nil {
		_ = console.WriteError(err.Error())
		return nil, err
	}
	return &execHookHandlerProvider{
		argsTemplate:         argsTemplate,
		filePathFilterRegexp: filePathFilter,
		console:              console,
	}, nil
}
