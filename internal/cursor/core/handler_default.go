package cursorcore

import "github.com/sviatsviatsviat/wat/internal/core"

// defaultHookHandler runs cmd with template bindings derived from hookData.
type defaultHookHandler struct {
	hookData HookDataCommon
}

// NewDefaultHookHandler returns a [core.HookHandler] that uses hookData for common placeholders only.
func NewDefaultHookHandler(hookData HookDataCommon) (core.HookHandler, error) {
	return defaultHookHandler{hookData: hookData}, nil
}

// Handle runs cmd with [core.TemplateBindings] from hookData and fixed stdout payload.
func (handler defaultHookHandler) Handle(cmd core.Command) core.HookHandlerResult {
	ctx := &core.HookContext{
		TemplateBindings: newTemplateBindingsCommon(handler.hookData),
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: DefaultHookResponseLine}
}
