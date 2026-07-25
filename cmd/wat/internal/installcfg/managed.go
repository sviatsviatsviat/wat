package installcfg

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ParseWatRunFlags extracts --agent and --event from a wat run shell command.
func ParseWatRunFlags(command string) (agent, event string, ok bool) {
	fields := strings.Fields(command)
	for i := 0; i < len(fields)-1; i++ {
		switch fields[i] {
		case "--agent":
			agent = fields[i+1]
		case "--event":
			event = fields[i+1]
		}
	}
	return agent, event, agent != "" && event != ""
}

// IsWatManagedCommand reports whether command is a wat-managed hook entry.
func IsWatManagedCommand(command, agent, event, watAbs string) bool {
	parsedAgent, parsedEvent, ok := ParseWatRunFlags(command)
	if !ok || parsedAgent != agent || parsedEvent != event {
		return false
	}
	return isWatRunExecutable(command, watAbs)
}

// IsWatManagedAgentCommand reports whether command is a wat-managed hook entry
// for agent, regardless of its native event.
func IsWatManagedAgentCommand(command, agent, watAbs string) bool {
	parsedAgent, _, ok := ParseWatRunFlags(command)
	if !ok || parsedAgent != agent {
		return false
	}
	return isWatRunExecutable(command, watAbs)
}

func isWatRunExecutable(command, watAbs string) bool {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) != 6 || fields[1] != "run" {
		return false
	}
	program := strings.Trim(fields[0], `"'`)
	watAbs = strings.TrimSpace(watAbs)
	if program == watAbs {
		return true
	}
	base := strings.TrimSuffix(strings.ToLower(filepath.Base(program)), ".exe")
	return base == "wat"
}

// WatRunCommand builds the shell command wat install writes for an agent event.
func WatRunCommand(watAbs, agent, event string) string {
	return fmt.Sprintf("%s run --agent %s --event %s", watAbs, agent, event)
}
