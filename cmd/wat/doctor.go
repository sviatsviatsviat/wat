package main

import (
	"flag"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/doctor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
)

func newDoctorCmd() *subcommandRunner {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	return &subcommandRunner{
		name:    "doctor",
		summary: "verify toolchain, script, cache, and install state",
		long: "Check Go toolchain, .wat/ project, build cache, and installed hook entries.\n\n" +
			"Missing install wiring and wat-on-PATH gaps are warnings (hooks simply will not run).\n" +
			"Exits 0 when all checks pass (warnings are allowed). Exits 4 when any check fails.",
		fs: fs,
		run: func() int {
			return runDoctor(defaultDoctorDeps())
		},
	}
}

type doctorDeps struct {
	doctor.Deps
}

func defaultDoctorDeps() doctorDeps {
	d := doctor.DefaultDeps()
	d.WatVersion = watModuleVersionFn()
	return doctorDeps{Deps: d}
}

func runDoctor(deps doctorDeps) int {
	proj := project.Deps{
		Getenv: deps.Getenv,
		Getwd:  deps.Getwd,
		Stat:   deps.Stat,
	}
	watDir, watErr := project.Resolve(proj)
	ctx := doctor.Context{WatDir: watDir, WatErr: watErr}
	results := doctor.Run(deps.Deps, ctx)

	failCount := doctor.FailCount(results)
	for _, r := range results {
		doctor.PrintResult(stdout, r)
	}
	if failCount > 0 {
		doctor.PrintFailureSummary(stdout, failCount)
		return exitCheckFailed
	}
	return exitOK
}
