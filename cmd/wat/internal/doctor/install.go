package doctor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
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

	results = append(results, checkInstalledCommands(agent, watAbs, expected, commands)...)
	results = append(results, checkExpectedCoverage(agent, watAbs, expected, commands)...)

	if len(results) == 0 && len(expected) > 0 {
		results = append(results, Result{
			Group:   "install",
			Status:  Pass,
			Message: fmt.Sprintf("%s: %d hook entries present", agent, len(expected)),
		})
	}
	return results
}

func checkInstalledCommands(agent, watAbs string, expected []string, commands []installedCommand) []Result {
	var results []Result
	for _, entry := range commands {
		if r, ok := validateWatCommand(entry.command); !ok {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: r,
				Fix:     "re-run wat install or fix the command in the config file",
			})
			continue
		}
		parsedAgent, event, _ := installcfg.ParseWatRunFlags(entry.command)
		if parsedAgent != agent {
			results = append(results, Result{
				Group:   "install",
				Status:  Warn,
				Message: fmt.Sprintf("%s: hooks[%q] has --agent %q", agent, entry.configEvent, parsedAgent),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		}
		if event != entry.configEvent {
			results = append(results, Result{
				Group:   "install",
				Status:  Warn,
				Message: fmt.Sprintf("%s: hooks[%q] has --event %q", agent, entry.configEvent, event),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		}
		if installcfg.IsWatManagedAgentCommand(entry.command, agent, watAbs) &&
			!slices.Contains(expected, event) {
			results = append(results, Result{
				Group:   "install",
				Status:  Fail,
				Message: fmt.Sprintf("%s: unregistered --event %q is still installed", agent, event),
				Fix:     fmt.Sprintf("run wat install --agent %s", agent),
			})
		}
	}
	return results
}

func checkExpectedCoverage(agent, watAbs string, expected []string, commands []installedCommand) []Result {
	var results []Result
	for _, event := range expected {
		count := 0
		for _, entry := range commands {
			if installcfg.IsWatManagedCommand(entry.command, agent, event, watAbs) {
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
