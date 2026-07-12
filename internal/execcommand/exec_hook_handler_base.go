package execcommand

import (
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// execHookHandlerBase holds exec template state shared by all exec hook handlers.
type execHookHandlerBase struct {
	argsTemplate []string
	console      cli.Console
}

// hookDefaultResponder writes the default hook protocol line for the active capability.
type hookDefaultResponder interface {
	WriteDefaultToHost()
}

// runExecWithBindings expands the argument template with bindings, runs the subprocess, writes the
// default hook protocol line via respond, and returns the handler result.
func (b execHookHandlerBase) runExecWithBindings(bindings templateBindings, respond hookDefaultResponder) core.HookHandlerResult {
	var code int
	renderedArgs, unknownPlaceholderKeys := renderTokens(b.argsTemplate, bindings)
	if len(unknownPlaceholderKeys) > 0 {
		_ = b.console.WriteErrorf("unknown template placeholders: %s\n", strings.Join(unknownPlaceholderKeys, ", "))
		code = cli.ExitBadInput
	} else {
		code = runSubprocess(b.console, renderedArgs)
	}
	respond.WriteDefaultToHost()
	return core.HookHandlerResult{Code: code}
}
