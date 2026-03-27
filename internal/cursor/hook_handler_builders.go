package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/core"
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

// cursorHookHandlerBuilders maps hook_event_name to a builder.
var cursorHookHandlerBuilders = map[string]cursorcore.HookHandlerBuilder{
	"afterShellExecution": cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterShellExecutionFields],
	"afterMCPExecution":   newDefaultHookHandlerBuilder,
	"afterFileEdit":       cursorcore.NewHookHandlerFromEventFields[cursorcore.AfterFileEditFields],
	"afterTabFileEdit":    newDefaultHookHandlerBuilder,
	"afterAgentResponse":  newDefaultHookHandlerBuilder,
	"afterAgentThought":   newDefaultHookHandlerBuilder,
	"sessionEnd":          newDefaultHookHandlerBuilder,
}

func newDefaultHookHandlerBuilder(_ []byte, hookData cursorcore.HookDataCommon) (core.HookHandler, error) {
	return cursorcore.NewDefaultHookHandler(hookData)
}
