package cursor

import "github.com/sviatsviatsviat/wat/internal/core"

// hookHandlerBuilder builds a [core.HookHandler] from parsed hook fields.
type hookHandlerBuilder func(hookData hookDataCommon) (core.HookHandler, error)

// cursorHookHandlerBuilders maps hook_event_name to a builder.
var cursorHookHandlerBuilders = map[string]hookHandlerBuilder{
	"afterShellExecution": newDefaultHookHandler,
	"afterMCPExecution":   newDefaultHookHandler,
	"afterFileEdit":       newDefaultHookHandler,
	"afterTabFileEdit":    newDefaultHookHandler,
	"afterAgentResponse":  newDefaultHookHandler,
	"afterAgentThought":   newDefaultHookHandler,
	"sessionEnd":          newDefaultHookHandler,
}
