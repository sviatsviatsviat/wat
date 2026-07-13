package main

import (
	"fmt"
	"os/exec"

	"github.com/sviatsviatsviat/wat/cmd/wat/checks"
)

type doctorDeps struct {
	checks.Deps
}

func defaultDoctorDeps() doctorDeps {
	return doctorDeps{Deps: checks.DefaultDeps()}
}

func (d doctorDeps) withProjectHooks(runDeps runDeps) doctorDeps {
	d.MustHaveWatFiles = func(watDir string) error {
		return mustHaveWatFiles(watDir, runDeps)
	}
	d.HookBuildCacheKey = func(watDir string) (string, error) {
		return hookBuildCacheKey(watDir, runDeps)
	}
	d.HooksBinaryPath = hooksBinaryPath
	return d
}

func runDoctor(deps doctorDeps) int {
	runDeps := deps.runDeps()
	deps = deps.withProjectHooks(runDeps)

	watDir, watErr := resolveWatDir(runDeps)
	ctx := checks.Context{WatDir: watDir, WatErr: watErr}
	results := checks.Run(deps.Deps, ctx)

	failCount := checks.FailCount(results)
	for _, r := range results {
		checks.PrintResult(stdout, r)
	}
	if failCount > 0 {
		_, _ = fmt.Fprintf(stdout, "\nwat doctor: %d check(s) failed\n", failCount)
		return exitCheckFailed
	}
	return exitOK
}

func (d doctorDeps) runDeps() runDeps {
	return runDeps{
		getenv:   d.Getenv,
		getwd:    d.Getwd,
		stat:     d.Stat,
		readDir:  d.ReadDir,
		readFile: d.ReadFile,
		mkdirAll: d.MkdirAll,
		command:  d.Command,
		runCmd: func(cmd *exec.Cmd) error {
			return cmd.Run()
		},
	}
}
