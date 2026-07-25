package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/paths"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// Install verifies wat on PATH and installed hook config entries.
// Missing agent configs, missing wat-managed entries, and wat not on PATH are
// warnings: they mean hooks will not be invoked, not that the project is broken.
func Install(deps Deps, ctx Context) []Result {
	var results []Result

	if ctx.ManifestErr != nil {
		return []Result{{
			Group:   "install",
			Status:  Fail,
			Message: "cannot verify install without authored hook registrations",
			Fix:     "fix .wat/hooks.go and re-run wat doctor",
		}}
	}
	watAbs, watPathErr := deps.LookPath("wat")
	if watPathErr != nil {
		results = append(results, Result{
			Group:   "install",
			Status:  Warn,
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
	configs := paths.All(root)

	present := make(map[string]os.FileInfo, len(configs))
	for _, cfg := range configs {
		fi, err := deps.Stat(cfg.Path)
		if err == nil {
			present[cfg.Path] = fi
		} else if !errors.Is(err, os.ErrNotExist) {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("stat %s failed", cfg.Path),
				Fix:     "fix permissions or re-run wat install",
			})
		}
	}
	for _, cfg := range configs {
		expected := ctx.Manifest.EventsFor(cfg.Agent)
		if _, ok := present[cfg.Path]; !ok {
			if len(expected) > 0 {
				results = append(results, Result{
					Group:   "install",
					Status:  Warn,
					Message: fmt.Sprintf("%s: hook config file missing", cfg.Agent),
					Fix:     fmt.Sprintf("run wat install --agent %s", cfg.Agent),
				})
			}
			continue
		}
		results = append(results, agentInstall(deps, cfg.Agent, cfg.Path, watAbs, expected)...)
	}

	results = append(results, disableAllHooks(deps, paths.ConfigPath(sdkclaude.Dialect, root))...)
	if len(ctx.Manifest.Registrations) == 0 && FailCount(results) == 0 {
		results = append(results, Result{
			Group:   "install",
			Status:  Pass,
			Message: "no hooks registered; no install entries required",
		})
	}

	return results
}

func agentInstall(deps Deps, agent, path, watAbs string, expected []string) []Result {
	var results []Result

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
			continue
		}
		_, event, _ := installcfg.ParseWatRunFlags(cmd)
		if installcfg.IsWatManagedAgentCommand(cmd, agent, watAbs) &&
			!slices.Contains(expected, event) {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("%s: unregistered --event %q is still installed", agent, event),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		}
	}

	for _, event := range expected {
		count := 0
		for _, cmd := range commands {
			if installcfg.IsWatManagedCommand(cmd, agent, event, watAbs) {
				count++
			}
		}
		switch count {
		case 0:
			results = append(results, Result{
				Group:   "install",
				Status:  Warn,
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

	if len(results) == 0 && len(expected) > 0 {
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
						return nil, fmt.Errorf("%s: malformed hook entry: %w", path, err)
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
	agent, _, ok := installcfg.ParseWatRunFlags(command)
	if !ok {
		return "invalid wat run command in config: " + firstLine(command), false
	}
	if err := dialect.Validate(agent); err != nil {
		return fmt.Sprintf("invalid --agent %q in config command", agent), false
	}
	return "", true
}
