package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// HookAdapterBuilder builds a [core.HookAdapter] from parsed hook fields.
type HookAdapterBuilder func(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error)

// cursorHookAdapterBuilders maps hook_event_name to a builder.
var cursorHookAdapterBuilders = map[string]HookAdapterBuilder{
	"afterShellExecution": hookAdapterFromEventFieldsAfterShellExecution,
	"afterMCPExecution":   hookAdapterFromEventFieldsAfterMCPExecution,
	"afterFileEdit":       hookAdapterFromEventFieldsAfterFileEdit,
	"afterTabFileEdit":    hookAdapterFromEventFieldsAfterFileEdit,
	"afterAgentResponse":  hookAdapterFromEventFieldsAfterAgentResponse,
	"afterAgentThought":   hookAdapterFromEventFieldsAfterAgentThought,
	"sessionEnd":          newDefaultHookAdapterBuilder,
}

func hookAdapterFromEventFieldsAfterShellExecution(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewHookAdapterFromEventFields[AfterShellExecutionFields](console, rawJSON, hookData)
}

func hookAdapterFromEventFieldsAfterMCPExecution(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewHookAdapterFromEventFields[AfterMCPExecutionFields](console, rawJSON, hookData)
}

func hookAdapterFromEventFieldsAfterFileEdit(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewHookAdapterFromEventFields[AfterFileEditFields](console, rawJSON, hookData)
}

func hookAdapterFromEventFieldsAfterAgentResponse(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewHookAdapterFromEventFields[AfterAgentResponseFields](console, rawJSON, hookData)
}

func hookAdapterFromEventFieldsAfterAgentThought(rawJSON []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewHookAdapterFromEventFields[AfterAgentThoughtFields](console, rawJSON, hookData)
}

func newDefaultHookAdapterBuilder(_ []byte, hookData HookDataCommon, console cli.Console) (core.HookAdapter, error) {
	return NewDefaultHookAdapter(console, hookData), nil
}

// NewHookAdapterFromEventFields parses rawJSON into [HookDataWithCommon] for T, then builds a
// hook adapter with [CursorHookRunData] (Common + EventSpecific from parsed fields).
func NewHookAdapterFromEventFields[T any](console cli.Console, rawJSON []byte, commonData HookDataCommon) (core.HookAdapter, error) {
	d, err := NewHookDataWithCommon[T](rawJSON, commonData)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor %s payload: %w", commonData.HookEventName, err)
	}
	return NewHookAdapter(console, d.HookDataCommon, &d.Fields), nil
}
