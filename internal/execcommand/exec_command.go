package execcommand

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

type execCommand struct {
	argsTemplate         []string
	filePathFilterRegexp *regexp.Regexp
	console              cli.Console
}

func (execCmd execCommand) Execute(hookContext *core.HookContext) int {
	if hookContext == nil {
		_ = execCmd.console.WriteError("internal error: HookContext is nil before Execute")
		return cli.ExitGeneral
	}
	if hookContext.HookHost != cursor.HookHostCursor {
		_ = execCmd.console.WriteError("internal error: exec command only supports Cursor hooks (unexpected HookHost)\n")
		return cli.ExitGeneral
	}

	parsed := hookContext.ParsedData
	if parsed == nil {
		_ = execCmd.console.WriteError("internal error: hook handler did not set HookContext.ParsedData before Execute\n")
		return cli.ExitGeneral
	}

	bindings, bindingsErr := templateBindingsForCursor(parsed)
	if bindingsErr != nil {
		_ = execCmd.console.WriteErrorf("internal error: %v\n", bindingsErr)
		return cli.ExitGeneral
	}

	if execCmd.filePathFilterRegexp != nil {
		if filePathFromHook, bindingDefined := bindings.TemplateValue("FILE_PATH"); bindingDefined {
			normalizedFilePath := filepath.ToSlash(filepath.Clean(filePathFromHook))
			if !execCmd.filePathFilterRegexp.MatchString(normalizedFilePath) {
				return cli.ExitSuccess
			}
		}
	}

	renderedArgs, unknownPlaceholderKeys := renderTokens(execCmd.argsTemplate, bindings)
	if len(unknownPlaceholderKeys) > 0 {
		_ = execCmd.console.WriteErrorf("unknown template placeholders: %s\n", strings.Join(unknownPlaceholderKeys, ", "))
		return cli.ExitBadInput
	}
	return runSubprocess(execCmd.console, renderedArgs)
}
