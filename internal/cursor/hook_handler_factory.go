// Package cursor implements Cursor hook stdin models (HookDataCommon, CursorHookRunData, event field types),
// generic CursorHookHandler and HookHandlerBuilder plumbing, and builds [core.HookHandler] values from stdin JSON.
package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// HookHandlerFactory builds Cursor hook handlers from stdin JSON.
type HookHandlerFactory struct{}

// NewHookHandlerFactory returns a factory for Cursor hook stdin and protocol.
func NewHookHandlerFactory() HookHandlerFactory {
	return HookHandlerFactory{}
}

func (f HookHandlerFactory) HookHandlerFromJSON(hookEventJSON []byte) (core.HookHandler, error) {
	if len(hookEventJSON) == 0 {
		return nil, fmt.Errorf("cursor hook stdin is empty or missing JSON object")
	}
	hookData, err := NewHookDataCommon(hookEventJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor hook JSON: %w", err)
	}
	builder, found := cursorHookHandlerBuilders[hookData.HookEventName]
	if !found {
		return nil, fmt.Errorf("cursor event %q is not supported yet", hookData.HookEventName)
	}
	return builder(hookEventJSON, hookData)
}
