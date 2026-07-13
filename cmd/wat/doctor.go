package main

import "flag"

func newDoctorCmd() *subcommandRunner {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	return &subcommandRunner{
		name:        "doctor",
		summary:     "verify toolchain, script, cache, and install state",
		long:        "Check Go toolchain, .wat/ project, build cache, and installed hook entries.",
		fs:          fs,
		notImplNote: "This command is not yet implemented.",
		run: func() int {
			return stubNotImplemented("doctor")
		},
	}
}
