package run

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/template"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

type runCommand struct {
	argvTemplate         []string
	filePathFilterRegexp *regexp.Regexp
	console              cli.Console
	subprocess           watexec.SubprocessRunner
}

func (runCmd runCommand) Execute(hookContext *core.HookContext) int {
	if hookContext == nil {
		_ = runCmd.console.WriteError("internal error: HookContext is nil before Execute")
		return cli.ExitGeneral
	}
	if hookContext.TemplateBindings == nil {
		_ = runCmd.console.WriteError("internal error: hook handler did not set HookContext.TemplateBindings before Execute")
		return cli.ExitGeneral
	}

	if runCmd.filePathFilterRegexp != nil {
		if filePathFromHook, bindingDefined := hookContext.TemplateBindings.TemplateValue("FILE_PATH"); bindingDefined {
			normalizedFilePath := filepath.ToSlash(filepath.Clean(filePathFromHook))
			if !runCmd.filePathFilterRegexp.MatchString(normalizedFilePath) {
				return cli.ExitSuccess
			}
		}
	}

	renderedArgv, unknownPlaceholderKeys := template.RenderTokens(runCmd.argvTemplate, hookContext.TemplateBindings)
	if len(unknownPlaceholderKeys) > 0 {
		_ = runCmd.console.WriteErrorf("unknown template placeholders: %s\n", strings.Join(unknownPlaceholderKeys, ", "))
		return cli.ExitBadInput
	}
	return runCmd.subprocess.Run(renderedArgv)
}
