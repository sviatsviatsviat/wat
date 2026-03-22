package cli

// Exit codes returned by wat and subprocess helpers.
const (
	ExitSuccess  = 0 // success
	ExitGeneral  = 1 // runtime or parse failure (not argv misuse)
	ExitBadInput = 2 // invalid argv, unknown command, template/host errors
)
