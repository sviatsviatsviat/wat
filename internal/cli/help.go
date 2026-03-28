package cli

// PrintRootHelp writes top-level usage text to console's diagnostic stream.
func PrintRootHelp(console Console) {
	_ = console.WriteError(rootHelpText)
}

// PrintExecHelp writes exec subcommand usage to console's diagnostic stream.
func PrintExecHelp(console Console) {
	_ = console.WriteError(execHelpText)
}
