package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// afterShellExecutionHookHandler runs cmd with template bindings from afterShellExecution hook payload.
type afterShellExecutionHookHandler struct {
	hookData hookDataAfterShellExecution
}

func newAfterShellExecutionHookHandler(rawJSON []byte, hookDataCommon hookDataCommon) (core.HookHandler, error) {
	hookData, err := newHookDataAfterShellExecution(rawJSON, hookDataCommon)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor afterShellExecution payload: %w", err)
	}
	return afterShellExecutionHookHandler{hookData: hookData}, nil
}

// Handle runs cmd with [core.TemplateBindings] from afterShellExecution payload and fixed stdout payload.
func (handler afterShellExecutionHookHandler) Handle(cmd core.Command) core.HookHandlerResult {
	ctx := &core.HookContext{
		TemplateBindings: newTemplateBindingsAfterShellExecution(handler.hookData),
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: defaultHookResponseLine}
}
