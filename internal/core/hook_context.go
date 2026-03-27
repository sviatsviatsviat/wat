package core

// HookContext carries host-specific parsed hook payload into [Command.Execute].
// ParsedData is set by the host hook handler; subcommands interpret it (e.g. [internal/run] for templating).
type HookContext struct {
	HookHost   string
	ParsedData any
}
