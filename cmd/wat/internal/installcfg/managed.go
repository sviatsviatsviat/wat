package installcfg

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	portclaude "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/claude"
	portcopilot "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/copilot"
	portcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
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
	switch dialect.Parse(agent) {
	case sdkclaude.Dialect:
		return eventInMapValues(portclaude.EventForKind, event)
	case sdkcopilot.Dialect:
		return eventInMapValues(portcopilot.EventForKind, event)
	case sdkcursor.Dialect:
		return eventInMapValues(portcursor.EventForKind, event) || portcursor.IsDedicatedEvent(event)
	default:
		return false
	}
}

// ExpectedEvents returns hook event names wat install writes for agent.
func ExpectedEvents(agent string) ([]string, error) {
	switch dialect.Parse(agent) {
	case sdkclaude.Dialect:
		return sortedValues(portclaude.EventForKind), nil
	case sdkcopilot.Dialect:
		return sortedValues(portcopilot.EventForKind), nil
	case sdkcursor.Dialect:
		eventSet := map[string]bool{}
		for _, ev := range sortedValues(portcursor.EventForKind) {
			eventSet[ev] = true
		}
		for ev := range portcursor.DedicatedEvents {
			eventSet[ev] = true
		}
		return sortedKeys(eventSet), nil
	default:
		return nil, fmt.Errorf("unknown agent %q", agent)
	}
}

func eventInMapValues(m map[model.Kind]string, event string) bool {
	for _, name := range m {
		if name == event {
			return true
		}
	}
	return false
}

func sortedValues[K comparable](m map[K]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		if v != "" {
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
