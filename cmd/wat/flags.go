package main

import (
	"flag"
	"fmt"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
)

// sharedFlags holds CLI flags reused across subcommands.
type sharedFlags struct {
	agent      *string
	event      *string
	failClosed *bool
}

func addAgentFlag(fs *flag.FlagSet, f *sharedFlags) {
	f.agent = fs.String("agent", "", "agent dialect dispatch hint (claude, copilot, cursor)")
}

func addEventFlag(fs *flag.FlagSet, f *sharedFlags) {
	f.event = fs.String("event", "", "native event dispatch hint (skips hook_event_name peek when set)")
}

func addFailClosedFlag(fs *flag.FlagSet, f *sharedFlags) {
	f.failClosed = fs.Bool("fail-closed", false, "treat build failures as block/deny (exit 2)")
}

func (f *sharedFlags) agentValue() string {
	if f == nil || f.agent == nil {
		return ""
	}
	return *f.agent
}

func (f *sharedFlags) eventValue() string {
	if f == nil || f.event == nil {
		return ""
	}
	return *f.event
}

func (f *sharedFlags) failClosedValue() bool {
	if f == nil || f.failClosed == nil {
		return false
	}
	return *f.failClosed
}

func validateAgent(agent string) error {
	return dialect.Validate(agent)
}

func validateSharedFlags(cmd string, f *sharedFlags) int {
	if f == nil || f.agent == nil {
		return exitOK
	}
	if err := validateAgent(*f.agent); err != nil {
		_, _ = fmt.Fprintf(stderr, "wat %s: %v\n", cmd, err)
		return exitUsage
	}
	return exitOK
}
