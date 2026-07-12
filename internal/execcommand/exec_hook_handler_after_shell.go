package execcommand

import "github.com/sviatsviatsviat/wat/internal/core"

// execHookHandlerAfterShell runs exec templates for afterShellExecution hooks.
type execHookHandlerAfterShell struct {
	execHookHandlerBase
	hook core.AfterShellExecutionHook
}

func (h execHookHandlerAfterShell) Handle() core.HookHandlerResult {
	bindings := templateBindingsAfterShellExecution{hook: h.hook}
	return h.runExecWithBindings(bindings, h.hook)
}
