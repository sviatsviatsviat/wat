// Package core defines host-neutral hook types: command execution, handlers, and template bindings.
package core

// Command runs once per hook invocation using placeholder values from ctx.
type Command interface {
	Execute(ctx *HookContext) int
}
