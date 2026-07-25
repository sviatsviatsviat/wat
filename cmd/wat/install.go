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
		summary: "install registered hooks into agent configs",
		long: "Build the .wat/ hook project, inspect its exported Hooks registrations, and reconcile hook configuration\n" +
			"entries for Claude Code, GitHub Copilot, and Cursor.\n\n" +
			"Only registered native events are installed. Stale wat-managed entries are removed.",
		fs: fs,
		run: func() int {
			plan, err := installcfg.ParseAgentPlan(*agent)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "wat install: %v\n", err)
				return exitUsage
			}
			if err := installcfg.Install(installcfg.Config{
				Agents:     plan,
				WatPath:    *watPath,
				WatVersion: watModuleVersionFn(),
			}, installcfg.DefaultDeps()); err != nil {
				_, _ = fmt.Fprintf(stderr, "wat install: %v\n", err)
				return exitRuntimeFailure
			}
			return exitOK
		},
	}
}
