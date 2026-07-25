package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hooktest"
)

func newTestCmd() *subcommandRunner {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	shared := &sharedFlags{}
	addAgentFlag(fs, shared)
	addEventFlag(fs, shared)
	fixture := fs.String("fixture", "", "path to fixture JSON, or \"-\" for stdin")
	verbose := fs.Bool("verbose", false, "print hook stderr")
	return &subcommandRunner{
		name:    "test",
		summary: "run hook script against fixture payloads",
		long: "Run the user's hook script against a fixture payload without invoking an agent.\n\n" +
			"wat test builds and executes the same cached .wat/hooks binary as wat run, feeding the\n" +
			"fixture on stdin. The fixture's native event must be registered. Pass --agent. It prints\n" +
			"fixture agent/event, the hook's stdout JSON, and exit code so you can iterate on handlers\n" +
			"locally.",
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

			cfg := hooktest.Config{
				Agent:   shared.agentValue(),
				Fixture: *fixture,
				Verbose: *verbose,
			}
			return hooktest.Run(cfg, watModuleVersionFn(), hooktest.DefaultDeps(stdout), os.Stdin, stderr)
		},
	}
}
