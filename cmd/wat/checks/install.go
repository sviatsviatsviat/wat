package checks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/internal/hookconfig"
)

// Install verifies wat on PATH and installed hook config entries.
func Install(deps Deps, ctx Context) []Result {
	var results []Result

	watAbs, watPathErr := deps.LookPath("wat")
	if watPathErr != nil {
		results = append(results, Result{
			Group:   "install",
			Status:  Fail,
			Message: "wat not found on PATH",
			Fix:     "install wat and ensure it is on PATH (or use wat install --wat-path)",
		})
		watAbs = "wat"
	} else {
		if abs, err := filepath.Abs(watAbs); err == nil {
			watAbs = abs
		}
		results = append(results, Result{
			Group:   "install",
			Status:  Pass,
			Message: "wat on PATH at " + watAbs,
		})
	}

	if ctx.WatErr != nil {
		results = append(results, Result{
			Group:   "install",
			Status:  Fail,
			Message: "cannot verify install without .wat/ project",
			Fix:     "run wat init",
		})
		return results
	}

	root := filepath.Dir(ctx.WatDir)

	type agentConfig struct {
		agent string
		path  string
	}
	configs := []agentConfig{
		{"claude", filepath.Join(root, ".claude", "settings.json")},
		{"copilot", filepath.Join(root, ".github", "hooks", "wat.json")},
		{"cursor", filepath.Join(root, ".cursor", "hooks.json")},
	}

	existing := 0
	for _, cfg := range configs {
		if _, err := deps.Stat(cfg.path); err == nil {
			existing++
		} else if !errors.Is(err, os.ErrNotExist) {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("stat %s failed", cfg.path),
				Fix:     "fix permissions or re-run wat install",
			})
		}
	}
	if existing == 0 {
		results = append(results, Result{
			Group:   "install",
			Status:  Fail,
			Message: "no hook config files found",
			Fix:     "run wat install",
		})
		return results
	}

	for _, cfg := range configs {
		if _, err := deps.Stat(cfg.path); err != nil {
			continue
		}
		results = append(results, agentInstall(deps, cfg.agent, cfg.path, watAbs)...)
	}

	results = append(results, disableAllHooks(deps, filepath.Join(root, ".claude", "settings.json"))...)

	return results
}

func agentInstall(deps Deps, agent, path, watAbs string) []Result {
	var results []Result
	expected, err := ExpectedInstallEvents(agent)
	if err != nil {
		return []Result{{
			Group:   "install",
			Status:  Fail,
			Message: err.Error(),
			Fix:     "re-run wat install",
		}}
	}

	commands, err := collectWatCommands(deps, agent, path)
	if err != nil {
		return []Result{{
			Group:   "install",
			Status:  Fail,
			Message: fmt.Sprintf("%s: %v", path, err),
			Fix:     fmt.Sprintf("run wat install --agent %s", agent),
		}}
	}

	for _, cmd := range commands {
		if r, ok := validateWatCommand(cmd); !ok {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: r,
				Fix:     "re-run wat install or fix the command in the config file",
			})
		}
	}

	for _, event := range expected {
		count := 0
		for _, cmd := range commands {
			if IsWatManagedCommand(cmd, agent, event, watAbs) {
				count++
			}
		}
		switch count {
		case 0:
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("%s: missing hook entry for %s", agent, event),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		case 1:
			continue
		default:
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("%s: duplicate wat-managed entries for %s", agent, event),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, Result{
			Group:   "install",
			Status:  Pass,
			Message: fmt.Sprintf("%s: %d hook entries present", agent, len(expected)),
		})
	}
	return results
}

func disableAllHooks(deps Deps, path string) []Result {
	if _, err := deps.Stat(path); err != nil {
		return nil
	}
	var settings claude.Settings
	if err := readInstallJSON(deps, path, &settings); err != nil {
		return nil
	}
	if settings.DisableAllHooks {
		return []Result{{
			Group:   "install",
			Status:  Warn,
			Message: "claude disableAllHooks is true",
			Fix:     "set disableAllHooks to false in .claude/settings.json",
		}}
	}
	return nil
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

func collectWatCommands(deps Deps, agent, path string) ([]string, error) {
	switch agent {
	case "claude":
		var settings claude.Settings
		if err := readInstallJSON(deps, path, &settings); err != nil {
			return nil, err
		}
		var cmds []string
		for _, groups := range settings.Hooks {
			for _, g := range groups {
				for _, raw := range g.Hooks {
					h, err := hookconfig.ParseHandler[claude.Handler](raw)
					if err != nil {
						continue
					}
					if h.Type == "" || h.Type == "command" {
						if strings.Contains(h.Command, "run --agent ") {
							cmds = append(cmds, h.Command)
						}
					}
				}
			}
		}
		return cmds, nil
	case "copilot":
		return collectFlatWatCommands[copilot.File](deps, path, copilot.ParseFlatCommand)
	case "cursor":
		return collectFlatWatCommands[cursor.File](deps, path, cursor.ParseFlatCommand)
	default:
		return nil, fmt.Errorf("unknown agent %q", agent)
	}
}

type flatHooksConfig interface {
	HooksMap() map[string][]json.RawMessage
}

func collectFlatWatCommands[F flatHooksConfig](deps Deps, path string, parseCommand func(json.RawMessage) (string, bool)) ([]string, error) {
	var f F
	if err := readInstallJSON(deps, path, &f); err != nil {
		return nil, err
	}
	return flatWatCommands(f.HooksMap(), parseCommand)
}

func flatWatCommands(hooks map[string][]json.RawMessage, getCommand func(json.RawMessage) (string, bool)) ([]string, error) {
	var cmds []string
	for _, handlers := range hooks {
		for _, raw := range handlers {
			c, ok := getCommand(raw)
			if ok && strings.Contains(c, "run --agent ") {
				cmds = append(cmds, c)
			}
		}
	}
	return cmds, nil
}

func validateWatCommand(command string) (string, bool) {
	agent, event, ok := ParseWatRunFlags(command)
	if !ok {
		return "invalid wat run command in config: " + firstLine(command), false
	}
	if err := validateAgent(agent); err != nil {
		return fmt.Sprintf("invalid --agent %q in config command", agent), false
	}
	if !isValidInstallEvent(agent, event) {
		return fmt.Sprintf("invalid --event %q for agent %s in config command", event, agent), false
	}
	return "", true
}
