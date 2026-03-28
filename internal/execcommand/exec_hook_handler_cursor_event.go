package execcommand

import (
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

// execHookHandlerCursorEvent runs exec templates for Cursor hooks whose flow is: build template
// bindings from hook data, run the subprocess, then write the default hook stdout line. Used for
// common-only payloads and for events such as afterShellExecution.
type execHookHandlerCursorEvent struct {
	execHookHandlerBase
	hook          core.HookAdapter
	buildBindings func() templateBindings
}

func (h execHookHandlerCursorEvent) Handle() core.HookHandlerResult {
	bindings := h.buildBindings()
	return h.runExecWithBindings(bindings, h.hook)
}

func execHookBindingsCommonOnly(common cursor.HookDataCommon) templateBindings {
	return newTemplateBindingsCommon(common)
}

func execHookBindingsAfterShellExecution(common cursor.HookDataCommon, event *cursor.AfterShellExecutionFields) templateBindings {
	return templateBindingsFromCursorEventPayload(common, event, afterShellExecutionPlaceholderExtractors)
}
