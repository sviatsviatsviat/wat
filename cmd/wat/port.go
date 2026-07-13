package main

import "flag"

func newPortCmd() *subcommandRunner {
	fs := flag.NewFlagSet("port", flag.ContinueOnError)
	return &subcommandRunner{
		name:        "port",
		summary:     "translate hook configs between agents",
		long:        "Convert hook configuration files between Claude Code, GitHub Copilot, and Cursor.",
		fs:          fs,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("port")
		},
	}
}
