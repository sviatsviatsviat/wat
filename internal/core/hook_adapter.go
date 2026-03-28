package core

// HookAdapter carries host-specific parsed hook payload for [HookHandlerProvider] and [HookHandler].
type HookAdapter interface {
	// HookHost is the hook host name (e.g. program first argument token) for this payload.
	HookHost() string
	// ReturnEmpty writes the default hook protocol response (e.g. Cursor "{}\n") using the console
	// captured when the adapter was built.
	ReturnEmpty()
}
