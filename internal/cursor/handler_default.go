package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/core"
)

// defaultHookResponseLine is the default hook stdout line (JSON object and newline).
const defaultHookResponseLine = "{}\n"

// defaultHookHandler runs cmd with template bindings derived from hookData.
type defaultHookHandler struct {
	hookData hookDataCommon
}

// newDefaultHookHandler returns a [core.HookHandler] that uses hookData for placeholders.
func newDefaultHookHandler(hookData hookDataCommon) (core.HookHandler, error) {
	return defaultHookHandler{hookData: hookData}, nil
}

// Handle runs cmd with [core.TemplateBindings] from hookData and fixed stdout payload.
func (handler defaultHookHandler) Handle(cmd core.Command) core.HookHandlerResult {
	ctx := &core.HookContext{
		TemplateBindings: newTemplateBindingsCommon(handler.hookData),
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: defaultHookResponseLine}
}
