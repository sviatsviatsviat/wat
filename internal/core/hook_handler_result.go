package core

// HookHandlerResult is the exit code and hook stdout payload from [HookHandler.Handle].
type HookHandlerResult struct {
	Code   int
	Output string
}
