package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/checks"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

type installAgentPlan struct {
	claude  bool
	copilot bool
	cursor  bool
}

func parseInstallAgent(s string) (installAgentPlan, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "", "all":
		return installAgentPlan{claude: true, copilot: true, cursor: true}, nil
	case "claude":
		return installAgentPlan{claude: true}, nil
	case "copilot":
		return installAgentPlan{copilot: true}, nil
	case "cursor":
		return installAgentPlan{cursor: true}, nil
	default:
		return installAgentPlan{}, fmt.Errorf("unknown agent %q (want claude, copilot, cursor, or all)", s)
	}
}

type installConfig struct {
	agents  installAgentPlan
	watPath string
}

type installDeps struct {
	getwd     func() (string, error)
	stat      func(string) (os.FileInfo, error)
	readFile  func(string) ([]byte, error)
	mkdirAll  func(string, os.FileMode) error
	writeFile func(string, []byte, os.FileMode) error
	lookPath  func(string) (string, error)
}

func defaultInstallDeps() installDeps {
	return installDeps{
		getwd:     os.Getwd,
		stat:      os.Stat,
		readFile:  os.ReadFile,
		mkdirAll:  os.MkdirAll,
		writeFile: os.WriteFile,
		lookPath:  exec.LookPath,
	}
}

func installProject(cfg installConfig, deps installDeps) error {
	wd, err := deps.getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	watDir, err := findWatDirForInstall(wd, deps)
	if err != nil {
		return err
	}
	root := filepath.Dir(watDir)

	watAbs, err := resolveWatPath(cfg.watPath, deps)
	if err != nil {
		return err
	}

	if cfg.agents.claude {
		if err := installClaude(root, watAbs, deps); err != nil {
			return err
		}
	}
	if cfg.agents.copilot {
		if err := installCopilot(root, watAbs, deps); err != nil {
			return err
		}
	}
	if cfg.agents.cursor {
		if err := installCursor(root, watAbs, deps); err != nil {
			return err
		}
	}

	return nil
}

func findWatDirForInstall(start string, deps installDeps) (string, error) {
	dir := start
	for {
		watDir := filepath.Join(dir, ".wat")
		if hasWatFiles(watDir, deps) {
			return watDir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .wat/ project found from %s (run \"wat init\" first)", start)
		}
		dir = parent
	}
}

func hasWatFiles(watDir string, deps installDeps) bool {
	hooksPath := filepath.Join(watDir, "hooks.go")
	if _, err := deps.stat(hooksPath); err != nil {
		return false
	}
	goModPath := filepath.Join(watDir, "go.mod")
	if _, err := deps.stat(goModPath); err != nil {
		return false
	}
	return true
}

func resolveWatPath(flagValue string, deps installDeps) (string, error) {
	if strings.TrimSpace(flagValue) != "" {
		abs, err := filepath.Abs(flagValue)
		if err == nil {
			return abs, nil
		}
		return flagValue, nil
	}
	p, err := deps.lookPath("wat")
	if err != nil {
		return "", fmt.Errorf("resolve wat executable: %w", err)
	}
	abs, err := filepath.Abs(p)
	if err == nil {
		return abs, nil
	}
	return p, nil
}

func installClaude(root, watAbs string, deps installDeps) error {
	path := filepath.Join(root, ".claude", "settings.json")
	settings := claude.Settings{Hooks: map[string][]claude.MatcherGroup{}}

	if err := readInstallJSON(path, &settings, deps); err != nil {
		return err
	}
	if settings.Hooks == nil {
		settings.Hooks = map[string][]claude.MatcherGroup{}
	}

	events, err := checks.ExpectedInstallEvents("claude")
	if err != nil {
		return err
	}
	for _, event := range events {
		cmd := watRunCommand(watAbs, "claude", event)
		settings.Hooks[event] = upsertClaudeGroups(settings.Hooks[event], cmd, "claude", event, watAbs)
	}

	return writeInstallJSON(path, settings, deps)
}

func upsertClaudeGroups(groups []claude.MatcherGroup, cmd, agent, event, watAbs string) []claude.MatcherGroup {
	// Remove wat-managed handlers from all groups.
	var out []claude.MatcherGroup
	for _, g := range groups {
		var kept []json.RawMessage
		for _, raw := range g.Hooks {
			h, err := hookkit.ParseHandler[claude.Handler](raw)
			if err != nil {
				kept = append(kept, raw)
				continue
			}
			if h.Type == "" || h.Type == "command" {
				if isWatManagedCommand(h.Command, agent, event, watAbs) {
					continue
				}
			}
			kept = append(kept, raw)
		}
		if len(kept) == 0 {
			continue
		}
		g.Hooks = kept
		out = append(out, g)
	}

	// Ensure a default matcher group exists and add the command handler.
	added := false
	for i := range out {
		if out[i].Matcher == "" && len(out[i].If) == 0 {
			raw, _ := hookkit.MarshalHandler(claude.Handler{Type: "command", Command: cmd})
			out[i].Hooks = append(out[i].Hooks, raw)
			added = true
			break
		}
	}
	if !added {
		raw, _ := hookkit.MarshalHandler(claude.Handler{Type: "command", Command: cmd})
		out = append(out, claude.MatcherGroup{Hooks: []json.RawMessage{raw}})
	}
	return out
}

func installCopilot(root, watAbs string, deps installDeps) error {
	path := filepath.Join(root, ".github", "hooks", "wat.json")
	f := copilot.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	if err := readInstallJSON(path, &f, deps); err != nil {
		return err
	}
	initFlatHooksFile(&f.Hooks, &f.Version)
	events, err := checks.ExpectedInstallEvents("copilot")
	if err != nil {
		return err
	}
	upsertFlatHookEvents(f.Hooks, "copilot", watAbs, events, copilot.ParseFlatCommand, marshalCopilotCommand)
	return writeInstallJSON(path, f, deps)
}

func installCursor(root, watAbs string, deps installDeps) error {
	path := filepath.Join(root, ".cursor", "hooks.json")
	f := cursor.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	if err := readInstallJSON(path, &f, deps); err != nil {
		return err
	}
	initFlatHooksFile(&f.Hooks, &f.Version)
	events, err := checks.ExpectedInstallEvents("cursor")
	if err != nil {
		return err
	}
	upsertFlatHookEvents(f.Hooks, "cursor", watAbs, events, cursor.ParseFlatCommand, marshalCursorCommand)
	return writeInstallJSON(path, f, deps)
}

func initFlatHooksFile(hooks *map[string][]json.RawMessage, version *int) {
	if *hooks == nil {
		*hooks = map[string][]json.RawMessage{}
	}
	if *version == 0 {
		*version = 1
	}
}

func upsertFlatHookEvents(hooks map[string][]json.RawMessage, agent, watAbs string, events []string, parseCommand func(json.RawMessage) (string, bool), marshalCommand func(string) json.RawMessage) {
	for _, event := range events {
		cmd := watRunCommand(watAbs, agent, event)
		hooks[event] = upsertFlatHandlers(hooks[event], cmd, agent, event, watAbs, parseCommand, marshalCommand)
	}
}

func marshalCopilotCommand(cmd string) json.RawMessage {
	raw, _ := copilot.MarshalFlatCommand(cmd)
	return raw
}

func marshalCursorCommand(cmd string) json.RawMessage {
	raw, _ := cursor.MarshalFlatCommand(cmd)
	return raw
}

func readInstallJSON(path string, dest any, deps installDeps) error {
	data, err := deps.readFile(path)
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

func writeInstallJSON(path string, v any, deps installDeps) error {
	if err := deps.mkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	out = append(out, '\n')
	if err := deps.writeFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func upsertFlatHandlers(existing []json.RawMessage, cmd, agent, event, watAbs string, getCommand func(json.RawMessage) (string, bool), makeHandler func(string) json.RawMessage) []json.RawMessage {
	var out []json.RawMessage
	for _, raw := range existing {
		c, ok := getCommand(raw)
		if ok && checks.IsWatManagedCommand(c, agent, event, watAbs) {
			continue
		}
		out = append(out, raw)
	}
	out = append(out, makeHandler(cmd))
	return out
}

func watRunCommand(watAbs, agent, event string) string {
	return fmt.Sprintf("%s run --agent %s --event %s", watAbs, agent, event)
}

func isWatManagedCommand(command, agent, event, watAbs string) bool {
	return checks.IsWatManagedCommand(command, agent, event, watAbs)
}
