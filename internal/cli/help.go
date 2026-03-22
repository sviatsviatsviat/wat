package cli

// PrintRootHelp writes top-level usage text to console's diagnostic stream.
func PrintRootHelp(console Console) {
	_ = console.WriteError(rootHelpText)
}

// PrintRunHelp writes run subcommand usage to console's diagnostic stream.
func PrintRunHelp(console Console) {
	_ = console.WriteError(runHelpText)
}
