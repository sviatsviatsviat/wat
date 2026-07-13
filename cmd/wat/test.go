package main

import "flag"

func newTestCmd() *subcommandRunner {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	addEventFlag(fs, shared)
	return &subcommandRunner{
		name:        "test",
		summary:     "run hook script against fixture payloads",
		long:        "Run the user's hook script against a fixture payload without invoking an agent.",
		fs:          fs,
		shared:      shared,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("test")
		},
	}
}
