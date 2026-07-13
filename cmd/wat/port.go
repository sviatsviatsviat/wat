package main

import (
	"flag"
	"fmt"
)

func newPortCmd() *subcommandRunner {
	fs := flag.NewFlagSet("port", flag.ContinueOnError)
	from := fs.String("from", "", "source agent dialect (claude, copilot, cursor)")
	to := fs.String("to", "", "target agent dialect (claude, copilot, cursor)")
	input := fs.String("input", "", "input config file (defaults to the source agent's well-known path)")
	output := fs.String("output", "", "output config file (defaults to stdout)")
	fs.StringVar(input, "i", "", "input config file (shorthand for -input)")
	fs.StringVar(output, "o", "", "output config file (shorthand for -output)")
	return &subcommandRunner{
		name:    "port",
		summary: "translate hook configs between agents",
		long: "Convert hook configuration files between Claude Code, GitHub Copilot, and Cursor.\n\n" +
			"Unlike `wat install`, which writes fresh wat-managed entries that invoke `wat run`, " +
			"`wat port` translates existing native hook configurations between agents. " +
			"Command strings are copied verbatim; shell script paths are not relocated or rewritten.\n\n" +
			"Lossy or unsupported mappings emit warnings on stderr but exit 0. Translation errors exit non-zero.",
		fs: fs,
		run: func() int {
			fromDialect, err := parsePortDialect(*from, "from")
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "wat port: %v\n", err)
				return exitUsage
			}
			toDialect, err := parsePortDialect(*to, "to")
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "wat port: %v\n", err)
				return exitUsage
			}
			return portProject(portConfig{
				from:       fromDialect,
				to:         toDialect,
				inputPath:  *input,
				outputPath: *output,
			}, defaultPortDeps())
		},
	}
}
