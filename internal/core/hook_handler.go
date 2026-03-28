package core

// HookHandler runs one hook invocation for a wat subcommand (e.g. exec).
type HookHandler interface {
	Handle() HookHandlerResult
}
