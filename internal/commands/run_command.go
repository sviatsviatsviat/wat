// Package commands implements wat subcommands as [core.Command] values.
package commands

import (
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/template"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

// NewRunCommand returns a [core.Command] that expands argvTemplate with hook placeholders and runs it via subprocessRunner.
// Empty argvTemplate writes diagnostics and run help and returns a non-nil error.
func NewRunCommand(console cli.Console, subprocessRunner watexec.SubprocessRunner, argvTemplate []string) (core.Command, error) {
	if len(argvTemplate) == 0 {
		_ = console.WriteError("missing command to run (argv for the subprocess, e.g. go version)")
		cli.PrintRunHelp(console)
		return nil, fmt.Errorf("missing command to run")
	}
	return runCommand{
		argvTemplate: argvTemplate,
		console:      console,
		subprocess:   subprocessRunner,
	}, nil
}

type runCommand struct {
	argvTemplate []string
	console      cli.Console
	subprocess   watexec.SubprocessRunner
}

func (cmd runCommand) Execute(hookCtx *core.HookContext) int {
	if hookCtx == nil {
		_ = cmd.console.WriteError("internal error: HookContext is nil before Execute")
		return cli.ExitGeneral
	}
	if hookCtx.TemplateBindings == nil {
		_ = cmd.console.WriteError("internal error: hook handler did not set HookContext.TemplateBindings before Execute")
		return cli.ExitGeneral
	}
	renderedArgv, unknownPlaceholders := template.RenderTokens(cmd.argvTemplate, hookCtx.TemplateBindings)
	if len(unknownPlaceholders) > 0 {
		_ = cmd.console.WriteErrorf("unknown template placeholders: %s\n", strings.Join(unknownPlaceholders, ", "))
		return cli.ExitBadInput
	}
	return cmd.subprocess.Run(renderedArgv)
}
