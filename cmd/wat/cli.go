package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	stdout = io.Writer(os.Stdout)
	stderr = io.Writer(os.Stderr)
)

type subcommand struct {
	name    string
	summary string
	newCmd  func() *subcommandRunner
}

type subcommandRunner struct {
	name        string
	summary     string
	long        string
	fs          *flag.FlagSet
	shared      *sharedFlags
	run         func() int
	notImplNote string
}

var subcommands = []subcommand{
	{name: "init", summary: "scaffold a .wat/ hook project", newCmd: newInitCmd},
	{name: "install", summary: "write hook config entries pointing at wat run", newCmd: newInstallCmd},
	{name: "run", summary: "execute .wat/hooks.go on hook invocation", newCmd: newRunCmd},
	{name: "port", summary: "translate hook configs between agents", newCmd: newPortCmd},
	{name: "test", summary: "run hook script against fixture payloads", newCmd: newTestCmd},
	{name: "doctor", summary: "verify toolchain, script, cache, and install state", newCmd: newDoctorCmd},
}

func run(args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(stderr, rootUsage())
		return exitUsage
	}

	switch args[0] {
	case "-h", "--help", "help":
		if len(args) == 1 {
			_, _ = fmt.Fprint(stdout, rootUsage())
			return exitOK
		}
		return subcommandHelp(args[1])
	}

	for i := range subcommands {
		if subcommands[i].name == args[0] {
			return subcommands[i].newCmd().execute(args[1:])
		}
	}

	_, _ = fmt.Fprintf(stderr, "wat: unknown command %q\n\n", args[0])
	_, _ = fmt.Fprint(stderr, rootUsage())
	return exitUsage
}

func subcommandHelp(name string) int {
	for i := range subcommands {
		if subcommands[i].name == name {
			_, _ = fmt.Fprint(stdout, subcommands[i].newCmd().helpText())
			return exitOK
		}
	}
	_, _ = fmt.Fprintf(stderr, "wat: unknown command %q\n\n", name)
	_, _ = fmt.Fprint(stderr, rootUsage())
	return exitUsage
}

func (c *subcommandRunner) execute(args []string) int {
	c.fs.SetOutput(io.Discard)
	c.fs.Usage = func() {}

	if err := c.fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprint(stdout, c.helpText())
			return exitOK
		}

		_, _ = fmt.Fprintf(stderr, "wat %s: %v\n\n", c.name, err)
		_, _ = fmt.Fprint(stderr, c.helpText())
		return exitUsage
	}

	if code := validateSharedFlags(c.name, c.shared); code != exitOK {
		return code
	}

	return c.run()
}

func (c *subcommandRunner) helpText() string {
	var b strings.Builder
	fmt.Fprintf(&b, "wat %s — %s\n\n", c.name, c.summary)
	fmt.Fprintf(&b, "Usage:\n\n\twat %s [flags]\n\n", c.name)
	if c.long != "" {
		fmt.Fprintf(&b, "%s\n\n", c.long)
	}
	if c.fs != nil {
		anyFlags := false
		c.fs.VisitAll(func(*flag.Flag) { anyFlags = true })

		if anyFlags {
			fmt.Fprint(&b, "Flags:\n\n")
			tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
			c.fs.VisitAll(func(f *flag.Flag) {
				_, _ = fmt.Fprintf(tw, "  -%s\t%s\n", f.Name, f.Usage)
			})
			_ = tw.Flush()
			fmt.Fprint(&b, "\n")
		}
	}
	if c.notImplNote != "" {
		fmt.Fprintf(&b, "%s\n", c.notImplNote)
	}
	return b.String()
}

func rootUsage() string {
	var b strings.Builder
	fmt.Fprint(&b, "wat — run and install Go hook scripts for coding agents\n\n")
	fmt.Fprint(&b, "Usage:\n\n\twat <command> [arguments]\n\n")
	fmt.Fprint(&b, "The commands are:\n\n")
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	for i := range subcommands {
		_, _ = fmt.Fprintf(tw, "\t%s\t%s\n", subcommands[i].name, subcommands[i].summary)
	}
	_ = tw.Flush()
	fmt.Fprint(&b, "\nUse \"wat <command> -h\" for more information about a command.\n")
	return b.String()
}
