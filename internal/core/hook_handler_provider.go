package core

import (
	"errors"
	"fmt"
)

// HookHandlerProvider builds a [HookHandler] for a given [HookAdapter] from subcommand configuration.
type HookHandlerProvider interface {
	HookHandlerFor(hook HookAdapter) (HookHandler, error)
}

// ErrHookAdapterNotSupported is the sentinel wrapped by [HookAdapterNotSupportedError].
var ErrHookAdapterNotSupported = errors.New("hook adapter is not supported for this subcommand")

// HookAdapterNotSupportedError returns an error wrapping [ErrHookAdapterNotSupported] with the
// concrete dynamic type of hook (or "<nil>") for diagnostics. Subcommands should use this instead
// of formatting the error ad hoc.
func HookAdapterNotSupportedError(hook HookAdapter) error {
	if hook == nil {
		return fmt.Errorf("%w: <nil>", ErrHookAdapterNotSupported)
	}
	return fmt.Errorf("%w: %T", ErrHookAdapterNotSupported, hook)
}
