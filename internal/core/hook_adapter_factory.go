package core

import "github.com/sviatsviatsviat/wat/internal/cli"

// HookAdapterFactory builds a [HookAdapter] from raw hook event JSON bytes.
type HookAdapterFactory interface {
	HookAdapterFromJSON(hookEventJSON []byte, console cli.Console) (HookAdapter, error)
}
