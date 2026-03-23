package cursor

import "github.com/sviatsviatsviat/wat/internal/core"

// hookHandlerBuilder builds a [core.HookHandler] from parsed hook fields.
type hookHandlerBuilder func(rawJSON []byte, hookData hookDataCommon) (core.HookHandler, error)

// cursorHookHandlerBuilders maps hook_event_name to a builder.
var cursorHookHandlerBuilders = map[string]hookHandlerBuilder{
	"afterShellExecution": newAfterShellExecutionHookHandler,
	"afterMCPExecution":   newDefaultHookHandlerBuilder,
	"afterFileEdit":       newDefaultHookHandlerBuilder,
	"afterTabFileEdit":    newDefaultHookHandlerBuilder,
	"afterAgentResponse":  newDefaultHookHandlerBuilder,
	"afterAgentThought":   newDefaultHookHandlerBuilder,
	"sessionEnd":          newDefaultHookHandlerBuilder,
}

func newDefaultHookHandlerBuilder(_ []byte, hookData hookDataCommon) (core.HookHandler, error) {
	return newDefaultHookHandler(hookData)
}
