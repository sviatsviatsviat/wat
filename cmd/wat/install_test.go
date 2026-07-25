package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookconfig"
	hostclaude "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/claude"
	hostcopilot "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/copilot"
	hostcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/hostconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	sdkrun "github.com/sviatsviatsviat/wat/sdk/run"
)

func TestInstallProject_freshInstallAll(t *testing.T) {
	dir := t.TempDir()
	project := filepath.Join(dir, "project")
	subdir := filepath.Join(project, "a", "b")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return subdir, nil }

	watAbs := filepath.Join(project, "bin", "wat")
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Claude: true, Copilot: true, Cursor: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	// Claude settings.json two-level shape.
	claudePath := filepath.Join(project, ".claude", "settings.json")
	claudeBytes, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("read %s: %v", claudePath, err)
	}
	var settings hostclaude.Settings
	if err := json.Unmarshal(claudeBytes, &settings); err != nil {
		t.Fatalf("parse %s: %v\n%s", claudePath, err, string(claudeBytes))
	}
	if len(settings.Hooks) == 0 {
		t.Fatalf("expected claude hooks to be populated")
	}
	foundClaude := false
	for event, groups := range settings.Hooks {
		for _, g := range groups {
			for _, raw := range g.Hooks {
				h := parseHandler[hostclaude.Handler](t, raw)
				if strings.Contains(h.Command, "run --agent claude --event "+event) {
					if !strings.HasPrefix(h.Command, watAbs+" ") {
						t.Fatalf("claude command should use watAbs prefix: %q", h.Command)
					}
					foundClaude = true
				}
			}
		}
	}
	if !foundClaude {
		t.Fatalf("expected at least one claude wat-managed hook command")
	}

	// Copilot hooks file flat shape.
	copilotPath := filepath.Join(project, ".github", "hooks", "wat.json")
	copilotBytes, err := os.ReadFile(copilotPath)
	if err != nil {
		t.Fatalf("read %s: %v", copilotPath, err)
	}
	var copilotFile hostcopilot.File
	if err := json.Unmarshal(copilotBytes, &copilotFile); err != nil {
		t.Fatalf("parse %s: %v\n%s", copilotPath, err, string(copilotBytes))
	}
	if copilotFile.Version != 1 {
		t.Fatalf("copilot version = %d, want 1", copilotFile.Version)
	}
	if len(copilotFile.Hooks) == 0 {
		t.Fatalf("expected copilot hooks to be populated")
	}

	// Cursor hooks file flat shape.
	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	cursorBytes, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatalf("read %s: %v", cursorPath, err)
	}
	var cursorFile hostcursor.File
	if err := json.Unmarshal(cursorBytes, &cursorFile); err != nil {
		t.Fatalf("parse %s: %v\n%s", cursorPath, err, string(cursorBytes))
	}
	if cursorFile.Version != 1 {
		t.Fatalf("cursor version = %d, want 1", cursorFile.Version)
	}
	if len(cursorFile.Hooks) == 0 {
		t.Fatalf("expected cursor hooks to be populated")
	}
}

func TestInstallProject_discoversAuthoredHooksIntegration(t *testing.T) {
	project := setupTestHookProject(t)
	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	watAbs := filepath.Join(project, "bin", "wat")

	if err := installcfg.Install(installcfg.Config{
		Agents:     installcfg.AgentPlan{Cursor: true},
		WatPath:    watAbs,
		WatVersion: watModuleVersionFn(),
	}, deps); err != nil {
		t.Fatal(err)
	}

	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	body, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatal(err)
	}
	var installed hostcursor.File
	if err := json.Unmarshal(body, &installed); err != nil {
		t.Fatal(err)
	}
	for _, event := range []string{
		cursor.EventBeforeShellExecution,
		cursor.EventAfterFileEdit,
		cursor.EventStop,
	} {
		if len(installed.Hooks[event]) != 1 {
			t.Fatalf("%s handlers = %d, want 1", event, len(installed.Hooks[event]))
		}
	}
	if _, ok := installed.Hooks[cursor.EventSessionStart]; ok {
		t.Fatalf("unregistered %s event was installed", cursor.EventSessionStart)
	}
}

func TestInstallProject_mergePreservesUnrelated(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Seed Cursor config with an unrelated hook entry.
	unrelated := hostcursor.File{
		Version: 1,
		Hooks: map[string][]json.RawMessage{
			"sessionStart": {
				mustHandlerRaw(t, hostcursor.Handler{Type: hostcursor.HandlerTypeCommand, Command: "echo unrelated"}),
			},
		},
	}
	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(cursorPath), 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(unrelated, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(cursorPath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	watAbs := filepath.Join(project, "bin", "wat")
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Cursor: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	gotBytes, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatalf("read %s: %v", cursorPath, err)
	}
	var got hostcursor.File
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("parse %s: %v", cursorPath, err)
	}
	raw := got.Hooks["sessionStart"][0]
	h := parseHandler[hostcursor.Handler](t, raw)
	if h.Command != "echo unrelated" {
		t.Fatalf("unrelated command changed: %q", h.Command)
	}
}

func TestInstallProject_reinstallReplacesWatEntries(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	watAbs := filepath.Join(project, "bin", "wat")
	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(cursorPath), 0o755); err != nil {
		t.Fatal(err)
	}
	seed := hostcursor.File{
		Version: 1,
		Hooks: map[string][]json.RawMessage{
			"preToolUse": {
				mustHandlerRaw(t, hostcursor.Handler{Type: hostcursor.HandlerTypeCommand, Command: watAbs + " run --agent cursor --event preToolUse"}),
			},
		},
	}
	b, _ := json.MarshalIndent(seed, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(cursorPath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Cursor: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	gotBytes, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatalf("read %s: %v", cursorPath, err)
	}
	var got hostcursor.File
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("parse %s: %v", cursorPath, err)
	}
	handlers := got.Hooks["preToolUse"]
	var watCount int
	for _, raw := range handlers {
		h := parseHandler[hostcursor.Handler](t, raw)
		if strings.Contains(h.Command, "run --agent cursor --event preToolUse") {
			watCount++
		}
	}
	if watCount != 1 {
		t.Fatalf("expected exactly 1 wat-managed preToolUse handler, got %d", watCount)
	}
}

func TestInstallProject_agentFiltering(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	watAbs := filepath.Join(project, "bin", "wat")
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Cursor: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	if _, err := os.Stat(filepath.Join(project, ".cursor", "hooks.json")); err != nil {
		t.Fatalf("cursor hooks not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".claude", "settings.json")); err == nil {
		t.Fatalf("claude settings.json should not be written when installing only cursor")
	}
	if _, err := os.Stat(filepath.Join(project, ".github", "hooks", "wat.json")); err == nil {
		t.Fatalf("copilot wat.json should not be written when installing only cursor")
	}
}

func TestInstallProject_emptyManifestDoesNotCreateConfigs(t *testing.T) {
	project := t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &sdkrun.Manifest{Version: 1}
	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Claude: true, Copilot: true, Cursor: true},
		WatPath:  filepath.Join(project, "bin", "wat"),
		Manifest: manifest,
	}, deps); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		filepath.Join(project, ".claude", "settings.json"),
		filepath.Join(project, ".github", "hooks", "wat.json"),
		filepath.Join(project, ".cursor", "hooks.json"),
	} {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("config %s was created for an empty manifest", path)
		}
	}
}

func TestUpsertClaudeGroups_preservesUnrelatedRemovesWatManaged(t *testing.T) {
	watAbs := "/bin/wat"
	event := claude.EventPreToolUse
	cmd := installcfg.WatRunCommand(watAbs, "claude", event)

	groups := []hostclaude.MatcherGroup{
		{
			Matcher: "Bash",
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo unrelated"}),
			},
		},
		{
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: watAbs + " run --agent claude --event " + event}),
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo keep-me"}),
			},
		},
	}

	got := installcfg.UpsertClaudeGroups(groups, cmd, "claude", event, watAbs)
	if len(got) != 2 {
		t.Fatalf("group count = %d, want 2", len(got))
	}
	if got[0].Matcher != "Bash" {
		t.Fatalf("first group matcher = %q, want Bash", got[0].Matcher)
	}
	if h := parseHandler[hostclaude.Handler](t, got[0].Hooks[0]); h.Command != "echo unrelated" {
		t.Fatalf("unrelated matcher-group command changed: %q", h.Command)
	}
	if h := parseHandler[hostclaude.Handler](t, got[1].Hooks[0]); h.Command != "echo keep-me" {
		t.Fatalf("unrelated default-group command removed: %q", h.Command)
	}
	watCount := countClaudeWatHandlers(t, got[1].Hooks, "claude", event)
	if watCount != 1 {
		t.Fatalf("default group wat handler count = %d, want 1", watCount)
	}
	if h := parseHandler[hostclaude.Handler](t, got[1].Hooks[len(got[1].Hooks)-1]); h.Command != cmd {
		t.Fatalf("new wat command = %q, want %q", h.Command, cmd)
	}
}

func TestUpsertClaudeGroups_dropsEmptyGroups(t *testing.T) {
	watAbs := "/bin/wat"
	event := claude.EventPreToolUse
	cmd := installcfg.WatRunCommand(watAbs, "claude", event)

	groups := []hostclaude.MatcherGroup{
		{
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: watAbs + " run --agent claude --event " + event}),
			},
		},
		{
			Matcher: "Write",
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo unrelated"}),
			},
		},
	}

	got := installcfg.UpsertClaudeGroups(groups, cmd, "claude", event, watAbs)
	if len(got) != 2 {
		t.Fatalf("group count = %d, want 2 (empty wat-only group removed, Write group kept, new default added)", len(got))
	}
	if got[0].Matcher != "Write" {
		t.Fatalf("first group matcher = %q, want Write", got[0].Matcher)
	}
	if got[1].Matcher != "" {
		t.Fatalf("second group should be default matcher group, got matcher %q", got[1].Matcher)
	}
	if countClaudeWatHandlers(t, got[1].Hooks, "claude", event) != 1 {
		t.Fatalf("expected exactly one wat handler in new default group")
	}
}

func TestUpsertClaudeGroups_multipleDefaultGroupsOnlyFirstGetsWat(t *testing.T) {
	watAbs := "/bin/wat"
	event := claude.EventPreToolUse
	cmd := installcfg.WatRunCommand(watAbs, "claude", event)

	groups := []hostclaude.MatcherGroup{
		{
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo first-default"}),
			},
		},
		{
			Hooks: []json.RawMessage{
				mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo second-default"}),
			},
		},
	}

	got := installcfg.UpsertClaudeGroups(groups, cmd, "claude", event, watAbs)
	if len(got) != 2 {
		t.Fatalf("group count = %d, want 2", len(got))
	}
	if h := parseHandler[hostclaude.Handler](t, got[0].Hooks[0]); h.Command != "echo first-default" {
		t.Fatalf("first default group unrelated hook changed: %q", h.Command)
	}
	if countClaudeWatHandlers(t, got[0].Hooks, "claude", event) != 1 {
		t.Fatalf("first default group should contain exactly one wat handler")
	}
	if countClaudeWatHandlers(t, got[1].Hooks, "claude", event) != 0 {
		t.Fatalf("second default group should not receive wat handler")
	}
	if h := parseHandler[hostclaude.Handler](t, got[1].Hooks[0]); h.Command != "echo second-default" {
		t.Fatalf("second default group structure changed: %q", h.Command)
	}
}

func TestInstallProject_claudeMergePreservesUnrelated(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	unrelated := hostclaude.Settings{
		Hooks: map[string][]hostclaude.MatcherGroup{
			claude.EventPreToolUse: {
				{
					Matcher: "Bash",
					Hooks: []json.RawMessage{
						mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: "echo unrelated"}),
					},
				},
			},
		},
	}
	claudePath := filepath.Join(project, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(claudePath), 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(unrelated, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(claudePath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	watAbs := filepath.Join(project, "bin", "wat")
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Claude: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	gotBytes, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("read %s: %v", claudePath, err)
	}
	var got hostclaude.Settings
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("parse %s: %v", claudePath, err)
	}
	groups := got.Hooks[claude.EventPreToolUse]
	if len(groups) < 2 {
		t.Fatalf("expected matcher group plus default wat group, got %d groups", len(groups))
	}
	if groups[0].Matcher != "Bash" {
		t.Fatalf("matcher group lost: %q", groups[0].Matcher)
	}
	if h := parseHandler[hostclaude.Handler](t, groups[0].Hooks[0]); h.Command != "echo unrelated" {
		t.Fatalf("unrelated matcher-group command changed: %q", h.Command)
	}
}

func TestInstallProject_claudeReinstallReplacesWatEntries(t *testing.T) {
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".wat"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, ".wat", "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	watAbs := filepath.Join(project, "bin", "wat")
	event := claude.EventPreToolUse
	seed := hostclaude.Settings{
		Hooks: map[string][]hostclaude.MatcherGroup{
			event: {
				{
					Hooks: []json.RawMessage{
						mustHandlerRaw(t, hostclaude.Handler{Type: "command", Command: watAbs + " run --agent claude --event " + event}),
					},
				},
			},
		},
	}
	claudePath := filepath.Join(project, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(claudePath), 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(seed, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(claudePath, b, 0o644); err != nil {
		t.Fatal(err)
	}

	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Claude: true},
		WatPath:  watAbs,
		Manifest: installTestManifest(),
	}, deps); err != nil {
		t.Fatalf("installProject: %v", err)
	}

	gotBytes, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatalf("read %s: %v", claudePath, err)
	}
	var got hostclaude.Settings
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("parse %s: %v", claudePath, err)
	}
	var watCount int
	for _, g := range got.Hooks[event] {
		for _, raw := range g.Hooks {
			h := parseHandler[hostclaude.Handler](t, raw)
			if strings.Contains(h.Command, "run --agent claude --event "+event) {
				watCount++
			}
		}
	}
	if watCount != 1 {
		t.Fatalf("expected exactly 1 wat-managed %s handler, got %d", event, watCount)
	}
}

func TestInstallProject_removesStaleManagedEvents(t *testing.T) {
	project := t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	watAbs := filepath.Join(project, "bin", "wat")
	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	if err := os.MkdirAll(filepath.Dir(cursorPath), 0o755); err != nil {
		t.Fatal(err)
	}
	seed := hostcursor.File{
		Version: 1,
		Hooks: map[string][]json.RawMessage{
			cursor.EventPostToolUse: {
				mustHandlerRaw(t, hostcursor.Handler{
					Type:    hostcursor.HandlerTypeCommand,
					Command: installcfg.WatRunCommand(filepath.Join(project, "old", "wat"), cursor.Dialect, cursor.EventPostToolUse),
				}),
				mustHandlerRaw(t, hostcursor.Handler{
					Type:    hostcursor.HandlerTypeCommand,
					Command: "echo keep-me",
				}),
			},
		},
	}
	body, err := json.MarshalIndent(seed, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cursorPath, append(body, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &sdkrun.Manifest{
		Version: 1,
		Registrations: []sdkrun.Registration{{
			Dialect:      cursor.Dialect,
			Event:        cursor.EventPreToolUse,
			HandlerCount: 1,
		}},
	}
	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Cursor: true},
		WatPath:  watAbs,
		Manifest: manifest,
	}, deps); err != nil {
		t.Fatal(err)
	}

	var got hostcursor.File
	installed, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(installed, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Hooks[cursor.EventPreToolUse]) != 1 {
		t.Fatalf("preToolUse handlers = %d, want 1", len(got.Hooks[cursor.EventPreToolUse]))
	}
	staleHandlers := got.Hooks[cursor.EventPostToolUse]
	if len(staleHandlers) != 1 {
		t.Fatalf("postToolUse handlers = %d, want unrelated handler only", len(staleHandlers))
	}
	if handler := parseHandler[hostcursor.Handler](t, staleHandlers[0]); handler.Command != "echo keep-me" {
		t.Fatalf("preserved command = %q, want echo keep-me", handler.Command)
	}
}

func TestInstallProject_removesStaleClaudeManagedEvents(t *testing.T) {
	project := t.TempDir()
	watDir := filepath.Join(project, ".wat")
	if err := os.MkdirAll(watDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "hooks.go"), []byte("package hooks\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(watDir, "go.mod"), []byte("module wat-hooks\ngo 1.26\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	watAbs := filepath.Join(project, "bin", "wat")
	claudePath := filepath.Join(project, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(claudePath), 0o755); err != nil {
		t.Fatal(err)
	}
	seed := hostclaude.Settings{
		Hooks: map[string][]hostclaude.MatcherGroup{
			claude.EventPostToolUse: {{
				Matcher: "Bash",
				Hooks: []json.RawMessage{
					mustHandlerRaw(t, hostclaude.Handler{
						Type:    "command",
						Command: installcfg.WatRunCommand(filepath.Join(project, "old", "wat"), claude.Dialect, claude.EventPostToolUse),
					}),
					mustHandlerRaw(t, hostclaude.Handler{
						Type:    "command",
						Command: "echo keep-claude",
					}),
				},
			}},
			claude.EventSessionEnd: {{
				Hooks: []json.RawMessage{
					mustHandlerRaw(t, hostclaude.Handler{
						Type:    "command",
						Command: installcfg.WatRunCommand(filepath.Join(project, "old", "wat"), claude.Dialect, claude.EventSessionEnd),
					}),
				},
			}},
		},
	}
	body, err := json.MarshalIndent(seed, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(claudePath, append(body, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &sdkrun.Manifest{
		Version: 1,
		Registrations: []sdkrun.Registration{{
			Dialect:      claude.Dialect,
			Event:        claude.EventPreToolUse,
			HandlerCount: 1,
		}},
	}
	deps := installcfg.DefaultDeps()
	deps.Getwd = func() (string, error) { return project, nil }
	if err := installcfg.Install(installcfg.Config{
		Agents:   installcfg.AgentPlan{Claude: true},
		WatPath:  watAbs,
		Manifest: manifest,
	}, deps); err != nil {
		t.Fatal(err)
	}

	var got hostclaude.Settings
	installed, err := os.ReadFile(claudePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(installed, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got.Hooks[claude.EventSessionEnd]; ok {
		t.Fatalf("expected empty SessionEnd event to be removed, got %#v", got.Hooks[claude.EventSessionEnd])
	}
	groups := got.Hooks[claude.EventPostToolUse]
	if len(groups) != 1 || len(groups[0].Hooks) != 1 {
		t.Fatalf("PostToolUse groups = %#v, want one group with unrelated handler only", groups)
	}
	if handler := parseHandler[hostclaude.Handler](t, groups[0].Hooks[0]); handler.Command != "echo keep-claude" {
		t.Fatalf("preserved command = %q, want echo keep-claude", handler.Command)
	}
	if len(got.Hooks[claude.EventPreToolUse]) != 1 {
		t.Fatalf("PreToolUse groups = %d, want 1", len(got.Hooks[claude.EventPreToolUse]))
	}
}

func installTestManifest() *sdkrun.Manifest {
	return &sdkrun.Manifest{
		Version: 1,
		Registrations: []sdkrun.Registration{
			{Dialect: claude.Dialect, Event: claude.EventPreToolUse, HandlerCount: 1},
			{Dialect: copilot.Dialect, Event: copilot.EventPreToolUse, HandlerCount: 1},
			{Dialect: cursor.Dialect, Event: cursor.EventPreToolUse, HandlerCount: 1},
		},
	}
}

func mustHandlerRaw[T any](t *testing.T, h T) json.RawMessage {
	t.Helper()
	raw, err := hookconfig.MarshalHandler(h)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func parseHandler[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()
	parsed, err := hookconfig.ParseHandler[T](raw)
	if err != nil {
		t.Fatalf("parse handler: %v", err)
	}
	return parsed
}

func countClaudeWatHandlers(t *testing.T, hooks []json.RawMessage, agent, event string) int {
	t.Helper()
	var n int
	for _, raw := range hooks {
		h := parseHandler[hostclaude.Handler](t, raw)
		if strings.Contains(h.Command, "run --agent "+agent+" --event "+event) {
			n++
		}
	}
	return n
}
