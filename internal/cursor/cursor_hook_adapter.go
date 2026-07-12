package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// cursorHook holds parsed Cursor hook stdin for subcommand handlers.
type cursorHook[T any] struct {
	CommonInput        HookDataCommon
	EventSpecificInput *T
	console            cli.Console
}

// DefaultCursorHook is the hook adapter for common-only stdin (no event payload).
type DefaultCursorHook = cursorHook[struct{}]

// AfterFileEditCursorHook is the hook adapter for afterFileEdit / afterTabFileEdit.
type AfterFileEditCursorHook = cursorHook[AfterFileEditFields]

// AfterShellExecutionCursorHook is the hook adapter for afterShellExecution.
type AfterShellExecutionCursorHook = cursorHook[AfterShellExecutionFields]

// AfterMCPExecutionCursorHook is the hook adapter for afterMCPExecution.
type AfterMCPExecutionCursorHook = cursorHook[AfterMCPExecutionFields]

// AfterAgentResponseCursorHook is the hook adapter for afterAgentResponse.
type AfterAgentResponseCursorHook = cursorHook[AfterAgentResponseFields]

// AfterAgentThoughtCursorHook is the hook adapter for afterAgentThought.
type AfterAgentThoughtCursorHook = cursorHook[AfterAgentThoughtFields]

// SessionEndCursorHook is the hook adapter for sessionEnd.
type SessionEndCursorHook = cursorHook[SessionEndFields]

// defaultHookResponseLine is the Cursor hook stdout line (JSON object and newline).
const defaultHookResponseLine = "{}\n"

// writeDefaultHookResponse writes the default Cursor hook protocol stdout line using the console captured at construction.
func (a *cursorHook[T]) writeDefaultHookResponse() {
	if a.console != nil {
		_ = a.console.Write(defaultHookResponseLine)
	}
}

// NewDefaultHookAdapter returns a [core.HookAdapter] with common fields only (no event-specific payload).
func NewDefaultHookAdapter(console cli.Console, hookData HookDataCommon) core.HookAdapter {
	return &cursorHook[struct{}]{
		CommonInput:        hookData,
		EventSpecificInput: nil,
		console:            console,
	}
}

// NewHookAdapter returns a [core.HookAdapter] with shared common fields and optional event-specific payload.
// For common-only hooks (T == struct{}), eventSpecific is nil.
func NewHookAdapter[T any](console cli.Console, common HookDataCommon, eventSpecific *T) core.HookAdapter {
	return &cursorHook[T]{
		CommonInput:        common,
		EventSpecificInput: eventSpecific,
		console:            console,
	}
}
