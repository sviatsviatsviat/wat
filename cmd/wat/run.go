package main

import "flag"

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
			cfg := runConfig{
				agent:      "",
				event:      "",
				failClosed: false,
			}
			if shared != nil {
				if shared.agent != nil {
					cfg.agent = *shared.agent
				}
				if shared.event != nil {
					cfg.event = *shared.event
				}
				if shared.failClosed != nil {
					cfg.failClosed = *shared.failClosed
				}
			}
			return runHook(cfg, defaultRunDeps())
		},
		fs:     fs,
		shared: shared,
	}
}
