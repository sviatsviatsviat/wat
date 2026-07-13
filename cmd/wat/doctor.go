package main

import "flag"

func newDoctorCmd() *subcommandRunner {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	return &subcommandRunner{
		name:    "doctor",
		summary: "verify toolchain, script, cache, and install state",
		long: "Check Go toolchain, .wat/ project, build cache, and installed hook entries.\n\n" +
			"Exits 0 when all checks pass (warnings are allowed). Exits 4 when any check fails.",
		fs: fs,
		run: func() int {
			return runDoctor(defaultDoctorDeps())
		},
	}
}
