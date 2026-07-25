package installcfg

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
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
	fields := strings.Fields(strings.TrimSpace(command))
	wantAbs := []string{strings.TrimSpace(watAbs), "run", "--agent", agent, "--event", event}
	wantPlain := []string{"wat", "run", "--agent", agent, "--event", event}
	return slices.Equal(fields, wantAbs) || slices.Equal(fields, wantPlain)
}

// WatRunCommand builds the shell command wat install writes for an agent event.
func WatRunCommand(watAbs, agent, event string) string {
	return fmt.Sprintf("%s run --agent %s --event %s", watAbs, agent, event)
}

// IsValidEvent reports whether event is a known install event for agent.
func IsValidEvent(agent, event string) bool {
	events, err := eventsForAgent(agent)
	if err != nil {
		return false
	}
	return slices.Contains(events, event)
}

// ExpectedEvents returns hook event names wat install writes for agent.
func ExpectedEvents(agent string) ([]string, error) {
	events, err := eventsForAgent(agent)
	if err != nil {
		return nil, err
	}
	out := slices.Clone(events)
	sort.Strings(out)
	return out, nil
}

func eventsForAgent(agent string) ([]string, error) {
	switch dialect.Parse(agent) {
	case sdkclaude.Dialect:
		return claudeEvents, nil
	case sdkcopilot.Dialect:
		return copilotEvents, nil
	case sdkcursor.Dialect:
		return cursorEvents, nil
	default:
		return nil, fmt.Errorf("unknown agent %q", agent)
	}
}
