package core

// HookHandlerResult is the process exit code from [HookHandler.Handle].
// Hook protocol output is written by [HookAdapter.ReturnEmpty] and is not carried here.
type HookHandlerResult struct {
	Code int
}
