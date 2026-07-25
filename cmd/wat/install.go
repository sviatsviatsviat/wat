package main

import (
	"flag"
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
)

func newInstallCmd() *subcommandRunner {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	agent := fs.String("agent", "all", "agent dialect to install for (claude, copilot, cursor, all)")
	watPath := fs.String("wat-path", "", "path to wat executable (defaults to resolving wat from PATH)")
	return &subcommandRunner{
		name:    "install",
		summary: "write hook config entries pointing at wat run",
		long: "Write or merge hook configuration entries for Claude Code, GitHub Copilot, and Cursor.\n\n" +
			"This command generates fresh wat-managed hook entries that invoke `wat run`. Use `wat port` to translate existing hook configurations between agents.",
		fs: fs,
		run: func() int {
			plan, err := installcfg.ParseAgentPlan(*agent)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "wat install: %v\n", err)
				return exitUsage
			}
			if err := installcfg.Install(installcfg.Config{
				Agents:  plan,
				WatPath: *watPath,
			}, installcfg.DefaultDeps()); err != nil {
				_, _ = fmt.Fprintf(stderr, "wat install: %v\n", err)
				return exitRuntimeFailure
			}
			return exitOK
		},
	}
}
