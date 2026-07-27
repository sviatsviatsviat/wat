package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
)

// installedCommand pairs a native config map key with its wat run command.
type installedCommand struct {
	configEvent string
	command     string
}

func readInstallJSON(deps Deps, path string, dest any) error {
	data, err := deps.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func collectWatCommands(deps Deps, agent, path string) ([]installedCommand, error) {
	switch agent {
	case "claude":
		return collectClaudeWatCommands(deps, path)
	case "copilot":
		return collectFlatWatCommands[copilot.File](deps, path, copilot.ParseFlatCommand)
	case "cursor":
		return collectFlatWatCommands[cursor.File](deps, path, cursor.ParseFlatCommand)
	default:
		return nil, fmt.Errorf("unknown agent %q", agent)
	}
}

func collectClaudeWatCommands(deps Deps, path string) ([]installedCommand, error) {
	var settings claude.Settings
	if err := readInstallJSON(deps, path, &settings); err != nil {
		return nil, err
	}
	var cmds []installedCommand
	for event, groups := range settings.Hooks {
		for _, g := range groups {
			for _, raw := range g.Hooks {
				h, err := hookconfig.ParseHandler[claude.Handler](raw)
				if err != nil {
					return nil, fmt.Errorf("%s: malformed hook entry: %w", path, err)
				}
				if h.Type == "" || h.Type == "command" {
					if strings.Contains(h.Command, "run --agent ") {
						cmds = append(cmds, installedCommand{configEvent: event, command: h.Command})
					}
				}
			}
		}
	}
	return cmds, nil
}

type flatHooksConfig interface {
	HooksMap() map[string][]json.RawMessage
}

func collectFlatWatCommands[F flatHooksConfig](deps Deps, path string, parseCommand func(json.RawMessage) (string, bool)) ([]installedCommand, error) {
	var f F
	if err := readInstallJSON(deps, path, &f); err != nil {
		return nil, err
	}
	return flatWatCommands(f.HooksMap(), parseCommand)
}

func flatWatCommands(hooks map[string][]json.RawMessage, getCommand func(json.RawMessage) (string, bool)) ([]installedCommand, error) {
	var cmds []installedCommand
	for event, handlers := range hooks {
		for _, raw := range handlers {
			c, ok := getCommand(raw)
			if ok && strings.Contains(c, "run --agent ") {
				cmds = append(cmds, installedCommand{configEvent: event, command: c})
			}
		}
	}
	return cmds, nil
}

func validateWatCommand(command string) (string, bool) {
	agent, _, ok := installcfg.ParseWatRunFlags(command)
	if !ok {
		return "invalid wat run command in config: " + firstLine(command), false
	}
	if err := dialect.Validate(agent); err != nil {
		return fmt.Sprintf("invalid --agent %q in config command", agent), false
	}
	return "", true
}
