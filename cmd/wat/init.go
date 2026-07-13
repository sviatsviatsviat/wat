package main

import "flag"

func newInitCmd() *subcommandRunner {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	return &subcommandRunner{
		name:        "init",
		summary:     "scaffold a .wat/ hook project",
		long:        "Create .wat/go.mod and .wat/hooks.go in the current working directory.",
		fs:          fs,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("init")
		},
	}
}
