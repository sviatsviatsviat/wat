package main

// CLI exit codes for the wat command.
const (
	exitOK             = 0
	exitUsage          = 1
	exitNotImplemented = 2
	exitRuntimeFailure = 3

	// exitBuildFailed is used when wat fails to compile the user hook script.
	// It intentionally matches exitUsage (1) because the calling agent treats
	// non-zero exits as a hook failure (fail-open by default).
	exitBuildFailed = 1

	// exitFailClosed is used when --fail-closed is set and the hook script build
	// fails; it intentionally matches exitNotImplemented (2) because many agents
	// interpret exit code 2 as block/deny.
	exitFailClosed = 2
)
