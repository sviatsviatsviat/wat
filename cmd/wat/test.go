package main

import (
	"flag"
	"fmt"
	"strings"
)

func newTestCmd() *subcommandRunner {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	addEventFlag(fs, shared)
	fixture := fs.String("fixture", "", "path to fixture JSON, or \"-\" for stdin")
	verbose := fs.Bool("verbose", false, "print expanded event fields and hook stderr")
	return &subcommandRunner{
		name:    "test",
		summary: "run hook script against fixture payloads",
		long: "Run the user's hook script against a fixture payload without invoking an agent.\n\n" +
			"wat test builds and executes the same cached .wat/hooks binary as wat run, feeding the\n" +
			"fixture on stdin. It prints a decoded unified event summary (via agenthooks codecs) and\n" +
			"the hook's stdout JSON plus exit code so you can iterate on handlers locally.",
		fs:     fs,
		shared: shared,
		run: func() int {
			if fs.NArg() > 0 {
				_, _ = fmt.Fprintln(stderr, "wat test: unexpected positional arguments")
				return exitUsage
			}
			if strings.TrimSpace(*fixture) == "" {
				_, _ = fmt.Fprintln(stderr, "wat test: --fixture is required")
				return exitUsage
			}

			cfg := testConfig{
				fixture: *fixture,
				verbose: *verbose,
			}
			if shared != nil {
				if shared.agent != nil {
					cfg.agent = *shared.agent
				}
				if shared.event != nil {
					cfg.event = *shared.event
				}
			}
			return runTest(cfg, defaultTestDeps())
		},
	}
}
