package app

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

// newHookHandlerFactory returns a [core.HookHandlerFactory] for host or a non-nil error if
// host is not supported.
func newHookHandlerFactory(host string) (core.HookHandlerFactory, error) {
	switch host {
	case "cursor":
		return cursor.NewHookHandlerFactory(), nil
	default:
		return nil, fmt.Errorf("host %q is not supported yet", host)
	}
}
