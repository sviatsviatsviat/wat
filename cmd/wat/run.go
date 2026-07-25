package main

import (
	"flag"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookrun"
)

func newRunCmd() *subcommandRunner {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	addEventFlag(fs, shared)
	addFailClosedFlag(fs, shared)
	return &subcommandRunner{
		name:    "run",
		summary: "execute .wat/hooks.go on hook invocation",
		long: "Build (if needed) and execute the user's hook script, passing stdin through untouched.\n\n" +
			"wat run is designed for low warm-path latency: it hashes .wat/hooks.go + .wat/go.mod + wat version to\n" +
			"compute a cache key, then executes a cached binary under .wat/.cache/ on subsequent invocations.",
		run: func() int {
			cfg := hookrun.Config{FailClosed: shared.failClosedValue()}
			return hookrun.Run(cfg, watModuleVersionFn(), hookrun.DefaultDeps(), stderr)
		},
		fs:     fs,
		shared: shared,
	}
}
