package main

import "flag"

func newRunCmd() *subcommandRunner {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	addEventFlag(fs, shared)
	addFailClosedFlag(fs, shared)
	return &subcommandRunner{
		name:        "run",
		summary:     "execute .wat/hooks.go on hook invocation",
		long:        "Build (if needed) and execute the user's hook script, passing stdin through untouched.",
		fs:          fs,
		shared:      shared,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("run")
		},
	}
}
