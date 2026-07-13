package main

import "flag"

func newInstallCmd() *subcommandRunner {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	return &subcommandRunner{
		name:        "install",
		summary:     "write hook config entries pointing at wat run",
		long:        "Write or merge hook configuration entries for Claude Code, GitHub Copilot, and Cursor.",
		fs:          fs,
		shared:      shared,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("install")
		},
	}
}
