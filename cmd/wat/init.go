package main

import (
	"flag"
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/initproj"
)

func newInitCmd() *subcommandRunner {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	force := fs.Bool("force", false, "overwrite existing .wat/hooks.go")
	return &subcommandRunner{
		name:    "init",
		summary: "scaffold a .wat/ hook project",
		long: "Create .wat/go.mod, .wat/hooks.go, and .wat/.gitignore under the optional path (default: current working directory).\n\n" +
			"Safe to re-run: existing go.mod, .gitignore, and testdata fixtures are preserved; .wat/hooks.go requires --force to overwrite.",
		fs: fs,
		run: func() int {
			args := fs.Args()
			if len(args) > 1 {
				_, _ = fmt.Fprintln(stderr, "wat init: expected at most one optional path argument")
				return exitUsage
			}
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			if err := initproj.Init(root, *force, watModuleVersionFn(), initproj.DefaultDeps(), stdout, stderr); err != nil {
				_, _ = fmt.Fprintf(stderr, "wat init: %v\n", err)
				return exitRuntimeFailure
			}
			return exitOK
		},
	}
}
