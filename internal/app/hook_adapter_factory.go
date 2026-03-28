package app

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func newHookAdapterFactory(host string) (core.HookAdapterFactory, error) {
	switch host {
	case "cursor":
		return cursor.NewHookAdapterFactory(), nil
	default:
		return nil, fmt.Errorf("host %q is not supported yet", host)
	}
}
