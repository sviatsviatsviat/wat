package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// HookHandlerBuilder builds a [core.HookHandler] from parsed hook fields.
type HookHandlerBuilder func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error)

// cursorHookHandlerBuilders maps hook_event_name to a builder.
var cursorHookHandlerBuilders = map[string]HookHandlerBuilder{
	"afterShellExecution": NewHookHandlerFromEventFields[AfterShellExecutionFields],
	"afterMCPExecution":   newDefaultHookHandlerBuilder,
	"afterFileEdit":       NewHookHandlerFromEventFields[AfterFileEditFields],
	"afterTabFileEdit":    newDefaultHookHandlerBuilder,
	"afterAgentResponse":  newDefaultHookHandlerBuilder,
	"afterAgentThought":   newDefaultHookHandlerBuilder,
	"sessionEnd":          newDefaultHookHandlerBuilder,
}

func newDefaultHookHandlerBuilder(_ []byte, hookData HookDataCommon) (core.HookHandler, error) {
	return NewDefaultHookHandler(hookData)
}

// NewHookHandlerFromEventFields parses rawJSON into [HookDataWithCommon] for T, then builds a
// [CursorHookHandler] with [CursorHookRunData] (Common + EventSpecific from parsed fields).
func NewHookHandlerFromEventFields[T any](rawJSON []byte, commonData HookDataCommon) (core.HookHandler, error) {
	d, err := NewHookDataWithCommon[T](rawJSON, commonData)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor %s payload: %w", commonData.HookEventName, err)
	}
	return NewHookHandler(CursorHookRunData[T]{
		Common:        d.HookDataCommon,
		EventSpecific: &d.Fields,
	}), nil
}
