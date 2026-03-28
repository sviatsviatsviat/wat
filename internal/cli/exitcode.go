package cli

// Exit codes returned by wat and subprocess helpers.
const (
	ExitSuccess  = 0 // success
	ExitGeneral  = 1 // runtime or parse failure
	ExitBadInput = 2 // invalid program arguments, unknown command
)
