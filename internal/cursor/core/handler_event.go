package cursorcore

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// DefaultHookResponseLine is the Cursor hook stdout line (JSON object and newline).
const DefaultHookResponseLine = "{}\n"

// HookHandlerBuilder builds a [core.HookHandler] from parsed hook fields.
type HookHandlerBuilder func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error)

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
