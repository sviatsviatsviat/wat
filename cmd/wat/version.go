package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/modver"
)

// watModuleVersionFn resolves the module/CLI version. Tests replace it.
var watModuleVersionFn = modver.Resolve

func newVersionCmd() *subcommandRunner {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	return &subcommandRunner{
		name:    "version",
		summary: "print the wat module / CLI version",
		long: "Print the same module version string used by wat init (go.mod pin) and the hook build cache key.\n\n" +
			"Tagged installs print the release tag (for example v0.2.0-alpha). Local builds with VCS\n" +
			"stamping print a pseudo-version. Builds without module or VCS version information fail.",
		fs: fs,
		run: func() int {
			return runVersion()
		},
	}
}

func runVersion() int {
	v := strings.TrimSpace(watModuleVersionFn())
	if v == "" {
		_, _ = fmt.Fprintln(stderr, "wat version: determine wat module version (build with -buildvcs=true or use a tagged build)")
		return exitRuntimeFailure
	}
	_, _ = fmt.Fprintln(stdout, v)
	return exitOK
}
