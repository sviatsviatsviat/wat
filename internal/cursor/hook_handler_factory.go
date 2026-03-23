// Package cursor builds [core.HookHandler] values from Cursor hook stdin JSON.
package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

// HookHandlerFactory constructs [core.HookHandler] instances for registered Cursor hook events.
type HookHandlerFactory struct{}

// NewHookHandlerFactory returns an empty [HookHandlerFactory].
func NewHookHandlerFactory() HookHandlerFactory {
	return HookHandlerFactory{}
}

// HookHandlerFromJSON implements [core.HookHandlerFactory].
// It returns an error when hookEventJSON is empty, the event name is unknown, or construction fails.
func (HookHandlerFactory) HookHandlerFromJSON(hookEventJSON []byte) (core.HookHandler, error) {
	if len(hookEventJSON) == 0 {
		return nil, fmt.Errorf("cursor hook stdin is empty or missing JSON object")
	}
	hookData, err := cursorcore.NewHookDataCommon(hookEventJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor hook JSON: %w", err)
	}
	builder, found := cursorHookHandlerBuilders[hookData.HookEventName]
	if !found {
		return nil, fmt.Errorf("cursor event %q is not supported yet", hookData.HookEventName)
	}
	return builder(hookEventJSON, hookData)
}
