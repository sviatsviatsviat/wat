// Package cursor implements Cursor hook stdin models (HookDataCommon, CursorHookRunData, event field types),
// hook adapter type aliases (e.g. [DefaultCursorHookAdapter]), HookAdapterBuilder plumbing, and builds [core.HookAdapter] values from stdin JSON.
package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// HookAdapterFactory builds Cursor hook adapters from stdin JSON.
type HookAdapterFactory struct{}

// NewHookAdapterFactory returns a factory for Cursor hook stdin and protocol.
func NewHookAdapterFactory() HookAdapterFactory {
	return HookAdapterFactory{}
}

func (f HookAdapterFactory) HookAdapterFromJSON(hookEventJSON []byte, console cli.Console) (core.HookAdapter, error) {
	if len(hookEventJSON) == 0 {
		return nil, fmt.Errorf("cursor hook stdin is empty or missing JSON object")
	}
	hookData, err := NewHookDataCommon(hookEventJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor hook JSON: %w", err)
	}
	builder, found := cursorHookAdapterBuilders[hookData.HookEventName]
	if !found {
		return nil, fmt.Errorf("cursor event %q is not supported yet", hookData.HookEventName)
	}
	return builder(hookEventJSON, hookData, console)
}
