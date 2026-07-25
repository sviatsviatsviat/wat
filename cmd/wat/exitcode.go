package main

// CLI exit codes for the wat command.
const (
	exitOK             = 0
	exitUsage          = 1
	exitRuntimeFailure = 3
	exitCheckFailed    = 4
)

// Build-failure exits for wat run/test are defined in hookrun/hooktest
// (ExitBuildFailed=1, ExitFailClosed=2) and returned directly from those packages.
