package main

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

// versionOverride, when non-empty, replaces the build-info module version.
// E2E sets it via -ldflags so linked git worktrees (which often omit VCS
// stamps) still exercise init/cache/version with a known pin.
var versionOverride string

func watModuleVersion() string {
	if v := strings.TrimSpace(versionOverride); v != "" {
		return v
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if v := strings.TrimSpace(info.Main.Version); v != "" && v != "(devel)" {
		return v
	}

	var revision string
	var t time.Time
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			parsed, err := time.Parse(time.RFC3339, s.Value)
			if err == nil {
				t = parsed.UTC()
			}
		}
	}
	if revision == "" || t.IsZero() {
		return ""
	}
	short := revision
	if len(short) > 12 {
		short = short[:12]
	}
	return fmt.Sprintf("v0.0.0-%s-%s", t.Format("20060102150405"), short)
}

var watModuleVersionFn = watModuleVersion

func newVersionCmd() *subcommandRunner {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	return &subcommandRunner{
		name:    "version",
		summary: "print the wat module / CLI version",
		long: "Print the same module version string used by wat init (go.mod pin) and the hook build cache key.\n\n" +
			"Tagged installs print the release tag (for example v0.1.1-alpha). Local builds with VCS\n" +
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
