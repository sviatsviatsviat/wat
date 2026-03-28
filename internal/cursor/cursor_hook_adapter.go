package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// cursorHookAdapter holds parsed Cursor hook stdin for subcommand handlers.
type cursorHookAdapter[T any] struct {
	CommonInput        HookDataCommon
	EventSpecificInput *T
	console            cli.Console
}

// DefaultCursorHookAdapter is the hook adapter for common-only stdin (no event payload).
type DefaultCursorHookAdapter = cursorHookAdapter[struct{}]

// AfterFileEditCursorHookAdapter is the hook adapter for afterFileEdit / afterTabFileEdit.
type AfterFileEditCursorHookAdapter = cursorHookAdapter[AfterFileEditFields]

// AfterShellExecutionCursorHookAdapter is the hook adapter for afterShellExecution.
type AfterShellExecutionCursorHookAdapter = cursorHookAdapter[AfterShellExecutionFields]

// defaultHookResponseLine is the Cursor hook stdout line (JSON object and newline).
const defaultHookResponseLine = "{}\n"

// HookHost returns [HookHostCursor].
func (a *cursorHookAdapter[T]) HookHost() string {
	return HookHostCursor
}

// ReturnEmpty writes the default Cursor hook protocol stdout line using the console captured at construction.
func (a *cursorHookAdapter[T]) ReturnEmpty() {
	if a.console != nil {
		_ = a.console.Write(defaultHookResponseLine)
	}
}

// NewDefaultHookAdapter returns a [core.HookAdapter] with common fields only (no event-specific payload).
func NewDefaultHookAdapter(console cli.Console, hookData HookDataCommon) (core.HookAdapter, error) {
	return &cursorHookAdapter[struct{}]{
		CommonInput:        hookData,
		EventSpecificInput: nil,
		console:            console,
	}, nil
}

// NewHookAdapter returns a [core.HookAdapter] with shared common fields and optional event-specific payload.
// For common-only hooks (T == struct{}), eventSpecific is nil.
func NewHookAdapter[T any](console cli.Console, common HookDataCommon, eventSpecific *T) core.HookAdapter {
	return &cursorHookAdapter[T]{
		CommonInput:        common,
		EventSpecificInput: eventSpecific,
		console:            console,
	}
}
