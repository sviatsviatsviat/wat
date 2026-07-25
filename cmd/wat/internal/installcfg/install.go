package installcfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/paths"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/project"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// AgentPlan selects which agent configs to install.
type AgentPlan struct {
	Claude  bool
	Copilot bool
	Cursor  bool
}

// ParseAgentPlan parses the --agent flag value for wat install.
func ParseAgentPlan(s string) (AgentPlan, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "", "all":
		return AgentPlan{Claude: true, Copilot: true, Cursor: true}, nil
	case "claude":
		return AgentPlan{Claude: true}, nil
	case "copilot":
		return AgentPlan{Copilot: true}, nil
	case "cursor":
		return AgentPlan{Cursor: true}, nil
	default:
		return AgentPlan{}, fmt.Errorf("unknown agent %q (want claude, copilot, cursor, or all)", s)
	}
}

// Config holds options for Install.
type Config struct {
	Agents  AgentPlan
	WatPath string
}

// Deps holds injectable dependencies for Install.
type Deps struct {
	Getwd     func() (string, error)
	Stat      func(string) (os.FileInfo, error)
	ReadFile  func(string) ([]byte, error)
	MkdirAll  func(string, os.FileMode) error
	WriteFile func(string, []byte, os.FileMode) error
	LookPath  func(string) (string, error)
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Getwd:     os.Getwd,
		Stat:      os.Stat,
		ReadFile:  os.ReadFile,
		MkdirAll:  os.MkdirAll,
		WriteFile: os.WriteFile,
		LookPath:  exec.LookPath,
	}
}

// Install writes or merges wat-managed hook config entries for the selected agents.
func Install(cfg Config, deps Deps) error {
	wd, err := deps.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}

	watDir, err := project.FindFrom(wd, project.Deps{Getwd: deps.Getwd, Stat: deps.Stat})
	if err != nil {
		return err
	}
	root := filepath.Dir(watDir)

	watAbs, err := resolveWatPath(cfg.WatPath, deps)
	if err != nil {
		return err
	}

	if cfg.Agents.Claude {
		if err := installClaude(root, watAbs, deps); err != nil {
			return err
		}
	}
	if cfg.Agents.Copilot {
		if err := installCopilot(root, watAbs, deps); err != nil {
			return err
		}
	}
	if cfg.Agents.Cursor {
		if err := installCursor(root, watAbs, deps); err != nil {
			return err
		}
	}

	return nil
}

func resolveWatPath(flagValue string, deps Deps) (string, error) {
	if strings.TrimSpace(flagValue) != "" {
		abs, err := filepath.Abs(flagValue)
		if err == nil {
			return abs, nil
		}
		return flagValue, nil
	}
	p, err := deps.LookPath("wat")
	if err != nil {
		return "", fmt.Errorf("resolve wat executable: %w", err)
	}
	abs, err := filepath.Abs(p)
	if err == nil {
		return abs, nil
	}
	return p, nil
}

func installClaude(root, watAbs string, deps Deps) error {
	path := paths.ConfigPath(sdkclaude.Dialect, root)
	settings := claude.Settings{Hooks: map[string][]claude.MatcherGroup{}}

	if err := readJSON(path, &settings, deps); err != nil {
		return err
	}
	if settings.Hooks == nil {
		settings.Hooks = map[string][]claude.MatcherGroup{}
	}

	events, err := ExpectedEvents(sdkclaude.Dialect)
	if err != nil {
		return err
	}
	for _, event := range events {
		cmd := WatRunCommand(watAbs, sdkclaude.Dialect, event)
		settings.Hooks[event] = UpsertClaudeGroups(settings.Hooks[event], cmd, sdkclaude.Dialect, event, watAbs)
	}

	return writeJSON(path, settings, deps)
}

// UpsertClaudeGroups removes wat-managed handlers and adds cmd to a default matcher group.
func UpsertClaudeGroups(groups []claude.MatcherGroup, cmd, agent, event, watAbs string) []claude.MatcherGroup {
	var out []claude.MatcherGroup
	for _, g := range groups {
		var kept []json.RawMessage
		for _, raw := range g.Hooks {
			h, err := hookconfig.ParseHandler[claude.Handler](raw)
			if err != nil {
				kept = append(kept, raw)
				continue
			}
			if h.Type == "" || h.Type == "command" {
				if IsWatManagedCommand(h.Command, agent, event, watAbs) {
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

	added := false
	for i := range out {
		if out[i].Matcher == "" && len(out[i].If) == 0 {
			raw, _ := hookconfig.MarshalHandler(claude.Handler{Type: "command", Command: cmd})
			out[i].Hooks = append(out[i].Hooks, raw)
			added = true
			break
		}
	}
	if !added {
		raw, _ := hookconfig.MarshalHandler(claude.Handler{Type: "command", Command: cmd})
		out = append(out, claude.MatcherGroup{Hooks: []json.RawMessage{raw}})
	}
	return out
}

func installCopilot(root, watAbs string, deps Deps) error {
	path := paths.ConfigPath(sdkcopilot.Dialect, root)
	f := copilot.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	if err := readJSON(path, &f, deps); err != nil {
		return err
	}
	initFlatHooksFile(&f.Hooks, &f.Version)
	events, err := ExpectedEvents(sdkcopilot.Dialect)
	if err != nil {
		return err
	}
	upsertFlatHookEvents(f.Hooks, sdkcopilot.Dialect, watAbs, events, copilot.ParseFlatCommand, marshalCopilotCommand)
	return writeJSON(path, f, deps)
}

func installCursor(root, watAbs string, deps Deps) error {
	path := paths.ConfigPath(sdkcursor.Dialect, root)
	f := cursor.File{Version: 1, Hooks: map[string][]json.RawMessage{}}
	if err := readJSON(path, &f, deps); err != nil {
		return err
	}
	initFlatHooksFile(&f.Hooks, &f.Version)
	events, err := ExpectedEvents(sdkcursor.Dialect)
	if err != nil {
		return err
	}
	upsertFlatHookEvents(f.Hooks, sdkcursor.Dialect, watAbs, events, cursor.ParseFlatCommand, marshalCursorCommand)
	return writeJSON(path, f, deps)
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
		cmd := WatRunCommand(watAbs, agent, event)
		hooks[event] = UpsertFlatHandlers(hooks[event], cmd, agent, event, watAbs, parseCommand, marshalCommand)
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

func readJSON(path string, dest any, deps Deps) error {
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

func writeJSON(path string, v any, deps Deps) error {
	if err := deps.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	out = append(out, '\n')
	if err := deps.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// UpsertFlatHandlers removes wat-managed handlers and appends a new command handler.
func UpsertFlatHandlers(existing []json.RawMessage, cmd, agent, event, watAbs string, getCommand func(json.RawMessage) (string, bool), makeHandler func(string) json.RawMessage) []json.RawMessage {
	var out []json.RawMessage
	for _, raw := range existing {
		c, ok := getCommand(raw)
		if ok && IsWatManagedCommand(c, agent, event, watAbs) {
			continue
		}
		out = append(out, raw)
	}
	out = append(out, makeHandler(cmd))
	return out
}
