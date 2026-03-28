package core

// HookContext carries host-specific parsed hook payload into [Command.Execute].
type HookContext struct {
	HookHost   string
	ParsedData any
}
