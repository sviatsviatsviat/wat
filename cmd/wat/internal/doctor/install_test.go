package doctor

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/installcfg"
)

func TestAgentInstall_results(t *testing.T) {
	watAbs := "/usr/bin/wat"
	claudePath := "/project/.claude/settings.json"

	t.Run("missing", func(t *testing.T) {
		deps := Deps{ReadFile: func(string) ([]byte, error) {
			return []byte(`{"hooks":{}}`), nil
		}}
		results := agentInstall(deps, "claude", claudePath, watAbs)
		if statusCount(results, Fail) == 0 {
			t.Fatalf("expected missing Fail results, got %#v", results)
		}
		if !hasMessageContaining(results, "missing hook entry") {
			t.Fatalf("expected missing message, got %#v", results)
		}
	})

	t.Run("single", func(t *testing.T) {
		body, err := claudeConfigJSON(watAbs, false)
		if err != nil {
			t.Fatal(err)
		}
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		results := agentInstall(deps, "claude", claudePath, watAbs)
		if len(results) != 1 || results[0].Status != Pass {
			t.Fatalf("want single Pass, got %#v", results)
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		body, err := claudeConfigJSON(watAbs, true)
		if err != nil {
			t.Fatal(err)
		}
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		results := agentInstall(deps, "claude", claudePath, watAbs)
		if !hasMessageContaining(results, "duplicate wat-managed entries") {
			t.Fatalf("expected duplicate message, got %#v", results)
		}
	})

	t.Run("invalid-command", func(t *testing.T) {
		body := []byte(`{
			"hooks": {
				"PreToolUse": [{
					"hooks": [{"type":"command","command":"/usr/bin/wat run --agent claude --event NotARealEvent"}]
				}]
			}
		}`)
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		results := agentInstall(deps, "claude", claudePath, watAbs)
		if !hasMessageContaining(results, "invalid --event") {
			t.Fatalf("expected invalid event message, got %#v", results)
		}
	})
}

func TestCollectWatCommands(t *testing.T) {
	t.Run("claude", func(t *testing.T) {
		body := []byte(`{
			"hooks": {
				"PreToolUse": [{
					"hooks": [
						{"type":"command","command":"/usr/bin/wat run --agent claude --event PreToolUse"},
						{"type":"prompt","prompt":"skip"},
						{"type":"command","command":"echo other"}
					]
				}]
			}
		}`)
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		cmds, err := collectWatCommands(deps, "claude", "settings.json")
		if err != nil {
			t.Fatal(err)
		}
		if len(cmds) != 1 || !strings.Contains(cmds[0], "PreToolUse") {
			t.Fatalf("cmds = %v", cmds)
		}
	})

	t.Run("claude malformed", func(t *testing.T) {
		body := []byte(`{"hooks":{"PreToolUse":[{"hooks":["not-an-object"]}]}}`)
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		_, err := collectWatCommands(deps, "claude", "settings.json")
		if err == nil || !strings.Contains(err.Error(), "malformed hook entry") {
			t.Fatalf("got err %v, want malformed hook entry", err)
		}
	})

	t.Run("copilot", func(t *testing.T) {
		body := []byte(`{
			"version": 1,
			"hooks": {
				"preToolUse": [
					{"type":"command","command":"/usr/bin/wat run --agent copilot --event preToolUse"},
					{"type":"prompt","prompt":"skip"}
				]
			}
		}`)
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		cmds, err := collectWatCommands(deps, "copilot", "wat.json")
		if err != nil {
			t.Fatal(err)
		}
		if len(cmds) != 1 {
			t.Fatalf("cmds = %v", cmds)
		}
	})

	t.Run("cursor", func(t *testing.T) {
		body := []byte(`{
			"version": 1,
			"hooks": {
				"preToolUse": [
					{"command":"/usr/bin/wat run --agent cursor --event preToolUse"},
					{"type":"prompt","prompt":"skip"}
				]
			}
		}`)
		deps := Deps{ReadFile: func(string) ([]byte, error) { return body, nil }}
		cmds, err := collectWatCommands(deps, "cursor", "hooks.json")
		if err != nil {
			t.Fatal(err)
		}
		if len(cmds) != 1 {
			t.Fatalf("cmds = %v", cmds)
		}
	})
}

func TestFlatWatCommands_filters(t *testing.T) {
	hooks := map[string][]json.RawMessage{
		"preToolUse": {
			json.RawMessage(`{"type":"command","command":"wat run --agent copilot --event preToolUse"}`),
			json.RawMessage(`{"type":"command","command":"echo hi"}`),
			json.RawMessage(`{"type":"prompt","prompt":"x"}`),
			json.RawMessage(`not-json`),
		},
	}
	cmds, err := flatWatCommands(hooks, func(raw json.RawMessage) (string, bool) {
		var h struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		}
		if json.Unmarshal(raw, &h) != nil {
			return "", false
		}
		if h.Type != "" && h.Type != "command" {
			return "", false
		}
		return h.Command, true
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cmds) != 1 || cmds[0] != "wat run --agent copilot --event preToolUse" {
		t.Fatalf("cmds = %v", cmds)
	}
}

func TestValidateWatCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		ok      bool
		wantSub string
	}{
		{name: "valid", command: "wat run --agent claude --event PreToolUse", ok: true},
		{name: "malformed flags", command: "wat run --agent claude", ok: false, wantSub: "invalid wat run command"},
		{name: "invalid agent", command: "wat run --agent nosuch --event PreToolUse", ok: false, wantSub: "invalid --agent"},
		{name: "invalid event", command: "wat run --agent claude --event NotReal", ok: false, wantSub: "invalid --event"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, ok := validateWatCommand(tt.command)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v (msg %q)", ok, tt.ok, msg)
			}
			if !tt.ok && !strings.Contains(msg, tt.wantSub) {
				t.Fatalf("msg = %q, want substring %q", msg, tt.wantSub)
			}
		})
	}
}

func TestInstall_reusesCachedStats(t *testing.T) {
	root := t.TempDir()
	watDir := root + "/.wat"
	claudePath := root + "/.claude/settings.json"
	if err := os.MkdirAll(root+"/.claude", 0o755); err != nil {
		t.Fatal(err)
	}
	body, err := claudeConfigJSON("/usr/bin/wat", false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(claudePath, body, 0o644); err != nil {
		t.Fatal(err)
	}

	statCalls := map[string]int{}
	deps := DefaultDeps()
	deps.LookPath = func(string) (string, error) { return "/usr/bin/wat", nil }
	baseStat := deps.Stat
	deps.Stat = func(path string) (os.FileInfo, error) {
		statCalls[path]++
		return baseStat(path)
	}
	deps.ReadFile = func(path string) ([]byte, error) {
		if path == claudePath {
			return body, nil
		}
		return nil, os.ErrNotExist
	}

	_ = Install(deps, Context{WatDir: watDir})
	// One Stat from the presence scan, one from disableAllHooks — not a third from the install loop.
	if n := statCalls[claudePath]; n != 2 {
		t.Fatalf("stat(%s) called %d times, want 2", claudePath, n)
	}
}

func claudeConfigJSON(watAbs string, duplicate bool) ([]byte, error) {
	events, err := installcfg.ExpectedEvents("claude")
	if err != nil {
		return nil, err
	}
	hooks := map[string][]map[string]any{}
	for _, event := range events {
		cmd := installcfg.WatRunCommand(watAbs, "claude", event)
		entries := []map[string]any{{
			"type":    "command",
			"command": cmd,
		}}
		if duplicate && event == events[0] {
			entries = append(entries, map[string]any{
				"type":    "command",
				"command": cmd,
			})
		}
		hooks[event] = []map[string]any{{
			"hooks": entries,
		}}
	}
	return json.Marshal(map[string]any{"hooks": hooks})
}

func statusCount(results []Result, status Status) int {
	n := 0
	for _, r := range results {
		if r.Status == status {
			n++
		}
	}
	return n
}

func hasMessageContaining(results []Result, sub string) bool {
	for _, r := range results {
		if strings.Contains(r.Message, sub) {
			return true
		}
	}
	return false
}
