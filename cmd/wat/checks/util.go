package checks

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
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

// ParseGoModDirective returns the go version from a go.mod file body.
func ParseGoModDirective(data []byte) (string, error) {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			v := strings.TrimSpace(strings.TrimPrefix(line, "go "))
			if v == "" {
				return "", fmt.Errorf("empty go directive")
			}
			if i := strings.IndexByte(v, ' '); i >= 0 {
				v = v[:i]
			}
			return v, nil
		}
	}
	return "", fmt.Errorf("no go directive found")
}

// ParseInstalledGoVersion extracts the semver from go version output.
func ParseInstalledGoVersion(output string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(output))
	for _, f := range fields {
		if strings.HasPrefix(f, "go") && len(f) > 2 {
			v := strings.TrimPrefix(f, "go")
			if v != "" {
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("no go version in output %q", strings.TrimSpace(output))
}

// GoVersionAtLeast reports whether installed satisfies required.
func GoVersionAtLeast(installed, required string) bool {
	inst := parseVersionParts(installed)
	req := parseVersionParts(required)
	for i := 0; i < len(req); i++ {
		iv := 0
		if i < len(inst) {
			iv = inst[i]
		}
		if iv > req[i] {
			return true
		}
		if iv < req[i] {
			return false
		}
	}
	return true
}

func parseVersionParts(v string) []int {
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			end := 0
			for end < len(p) && p[end] >= '0' && p[end] <= '9' {
				end++
			}
			if end == 0 {
				continue
			}
			n, err = strconv.Atoi(p[:end])
			if err != nil {
				continue
			}
		}
		out = append(out, n)
	}
	return out
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func validateAgent(agent string) error {
	if agent == "" {
		return nil
	}
	if dialect.Parse(agent) == "" {
		return fmt.Errorf("unknown agent dialect %q (want claude, copilot, or cursor)", agent)
	}
	return nil
}

func isValidInstallEvent(agent, event string) bool {
	switch dialect.Parse(agent) {
	case sdkclaude.Dialect:
		return eventInMapValues(portclaude.EventForKind, event)
	case sdkcopilot.Dialect:
		canonical, ok := sdkcopilot.CanonicalEventName(event)
		if !ok {
			return false
		}
		return eventInMapValues(portcopilot.EventForKind, canonical)
	case sdkcursor.Dialect:
		return eventInMapValues(portcursor.EventForKind, event) || portcursor.IsDedicatedEvent(event)
	default:
		return false
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

// ExpectedInstallEvents returns hook event names wat install writes for agent.
func ExpectedInstallEvents(agent string) ([]string, error) {
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
