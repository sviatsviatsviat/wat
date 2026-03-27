package cursorcore

import "github.com/sviatsviatsviat/wat/internal/core"

// CursorHookHandler runs cmd with parsed hook data for subcommands.
type CursorHookHandler[T any] struct {
	runData CursorHookRunData[T]
}

// NewDefaultHookHandler returns a [core.HookHandler] with common fields only (no event-specific payload).
func NewDefaultHookHandler(hookData HookDataCommon) (core.HookHandler, error) {
	return CursorHookHandler[struct{}]{
		runData: CursorHookRunData[struct{}]{
			Common:        hookData,
			EventSpecific: nil,
		},
	}, nil
}

// NewHookHandler returns a [core.HookHandler] that passes runData into [core.HookContext.ParsedData].
func NewHookHandler[T any](runData CursorHookRunData[T]) core.HookHandler {
	return CursorHookHandler[T]{runData: runData}
}

// Handle runs cmd with [core.HookContext] and fixed stdout payload.
func (handler CursorHookHandler[T]) Handle(cmd core.Command) core.HookHandlerResult {
	rd := handler.runData
	ctx := &core.HookContext{
		HookHost:   HookHostCursor,
		ParsedData: &rd,
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: DefaultHookResponseLine}
}
