package core

// HookHandlerFactory builds a [HookHandler] from raw hook event JSON bytes.
type HookHandlerFactory interface {
	HookHandlerFromJSON(hookEventJSON []byte) (HookHandler, error)
}
