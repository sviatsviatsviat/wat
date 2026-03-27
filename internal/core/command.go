// Package core defines host-neutral hook types: command execution, handlers, and hook context.
package core

// Command runs once per hook invocation using ctx.
type Command interface {
	Execute(ctx *HookContext) int
}
