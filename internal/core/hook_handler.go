package core

// HookHandler runs cmd for one hook event and returns exit code and hook protocol output.
type HookHandler interface {
	Handle(cmd Command) HookHandlerResult
}
