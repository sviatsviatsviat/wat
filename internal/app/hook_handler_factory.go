package app

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func newHookHandlerFactory(execCtx core.WatExecutionContext) (core.HookHandlerFactory, error) {
	switch execCtx.Host() {
	case "cursor":
		return cursor.NewHookHandlerFactory(execCtx), nil
	default:
		return nil, fmt.Errorf("host %q is not supported yet", execCtx.Host())
	}
}
