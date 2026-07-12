package portconfig

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/agenthooks"
)

const claudeSettings = `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {"type": "command", "command": ".claude/hooks/block-rm.sh", "timeout": 15}
        ]
      },
      {
        "matcher": "Edit|Write",
        "hooks": [
          {"type": "command", "command": ".claude/hooks/lint.sh"}
        ]
      }
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": ".claude/hooks/require-tests.sh"}]}
    ],
    "MessageDisplay": [
      {"hooks": [{"type": "command", "command": ".claude/hooks/plain.sh"}]}
    ]
  }
}`

const copilotSettings = `{
  "version": 1,
  "hooks": {
    "preToolUse": [
      {"type": "command", "command": ".claude/hooks/block-rm.sh", "matcher": "bash", "timeoutSec": 15},
      {"type": "command", "command": ".claude/hooks/lint.sh", "matcher": "edit|write"}
    ],
    "agentStop": [
      {"type": "command", "command": ".claude/hooks/require-tests.sh"}
    ]
  }
}`

const cursorSettings = `{
  "version": 1,
  "hooks": {
    "beforeSubmitPrompt": [{"command": ".cursor/hooks/audit.sh"}],
    "preToolUse": [{"command": ".cursor/hooks/guard.sh", "matcher": "Shell"}]
  }
}`

const cursorDedicatedSettings = `{
  "version": 1,
  "hooks": {
    "beforeShellExecution": [{"command": ".cursor/hooks/block-force-push.sh", "timeout": 20}]
  }
}`

func TestParseClaude_mappableAndExtras(t *testing.T) {
	cfg, _, err := Parse([]byte(claudeSettings), agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[agenthooks.KindPreTool]
	if len(pre) != 2 {
		t.Fatalf("want 2 PreToolUse entries, got %d", len(pre))
	}
	if pre[0].Command != ".claude/hooks/block-rm.sh" || pre[0].Matcher != "Bash" || pre[0].TimeoutSec != 15 {
		t.Fatalf("first PreToolUse entry: %+v", pre[0])
	}
	if len(cfg.Hooks[agenthooks.KindStop]) != 1 {
		t.Fatalf("want 1 Stop entry, got %d", len(cfg.Hooks[agenthooks.KindStop]))
	}
	if len(cfg.Extras) != 1 || cfg.Extras[0].Event != "MessageDisplay" {
		t.Fatalf("MessageDisplay should be in Extras: %+v", cfg.Extras)
	}
}

func TestParseCopilot_timeoutAndCommand(t *testing.T) {
	cfg, _, err := Parse([]byte(copilotSettings), agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[agenthooks.KindPreTool]
	if len(pre) != 2 {
		t.Fatalf("want 2 preToolUse entries, got %d", len(pre))
	}
	foundTimeout := false
	for _, e := range pre {
		if e.Matcher == "bash" {
			if e.TimeoutSec != 15 {
				t.Errorf("timeoutSec = %d, want 15", e.TimeoutSec)
			}
			foundTimeout = true
		}
	}
	if !foundTimeout {
		t.Fatal("bash matcher entry not found")
	}
}

func TestParseCopilot_timeoutAlias(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"command","command":"x.sh","timeout":42}]}}`
	cfg, _, err := Parse([]byte(raw), agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Hooks[agenthooks.KindPreTool][0].TimeoutSec != 42 {
		t.Fatalf("timeout alias not resolved: %+v", cfg.Hooks[agenthooks.KindPreTool][0])
	}
}

func TestParseCursor_dedicatedEvent(t *testing.T) {
	cfg, _, err := Parse([]byte(cursorDedicatedSettings), agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[agenthooks.KindPreTool]
	if len(pre) != 1 {
		t.Fatalf("want 1 entry, got %d", len(pre))
	}
	if pre[0].NativeEvent != "beforeShellExecution" {
		t.Fatalf("NativeEvent = %q, want beforeShellExecution", pre[0].NativeEvent)
	}
	if pre[0].Command != ".cursor/hooks/block-force-push.sh" || pre[0].TimeoutSec != 20 {
		t.Fatalf("entry: %+v", pre[0])
	}
}

func TestEmitClaude_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(claudeSettings), agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg1.Hooks[agenthooks.KindPreTool]
	if len(pre) == 0 {
		t.Fatal("expected PreTool entries")
	}
	pre[0].TimeoutSec = 99
	out, _, err := Emit(cfg1, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Hooks[agenthooks.KindPreTool][0].TimeoutSec != 99 {
		t.Fatalf("mutated timeout not emitted: %+v", cfg2.Hooks[agenthooks.KindPreTool][0])
	}
	cfg1.Hooks[agenthooks.KindPreTool][0].Raw = cfg2.Hooks[agenthooks.KindPreTool][0].Raw
	assertClaudeConfigEqual(t, cfg1, cfg2)
}

func TestEmitCopilot_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(copilotSettings), agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg1.Hooks[agenthooks.KindPreTool]
	if len(pre) == 0 {
		t.Fatal("expected PreTool entries")
	}
	for i := range pre {
		if pre[i].Matcher == "bash" {
			pre[i].TimeoutSec = 99
			break
		}
	}
	out, _, err := Emit(cfg1, agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out, agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range cfg2.Hooks[agenthooks.KindPreTool] {
		if e.Matcher == "bash" && e.TimeoutSec != 99 {
			t.Fatalf("mutated timeout not emitted: %+v", e)
		}
	}
	for i := range cfg1.Hooks[agenthooks.KindPreTool] {
		for j := range cfg2.Hooks[agenthooks.KindPreTool] {
			if cfg1.Hooks[agenthooks.KindPreTool][i].Matcher == cfg2.Hooks[agenthooks.KindPreTool][j].Matcher {
				cfg1.Hooks[agenthooks.KindPreTool][i].Raw = cfg2.Hooks[agenthooks.KindPreTool][j].Raw
			}
		}
	}
	assertCopilotHooksEqual(t, cfg1, cfg2)
}

func TestEmitCursor_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(cursorSettings), agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	assertCursorHooksEqual(t, cfg1, cfg2)
}

func TestEmitCursor_dedicatedEventRoundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(cursorDedicatedSettings), agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `"beforeShellExecution"`) {
		t.Errorf("beforeShellExecution not preserved: %s", s)
	}
	cfg2, _, err := Parse(out, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg2.Hooks[agenthooks.KindPreTool]) != 1 {
		t.Fatalf("want 1 pre-tool entry after round-trip, got %d", len(cfg2.Hooks[agenthooks.KindPreTool]))
	}
	e := cfg2.Hooks[agenthooks.KindPreTool][0]
	if e.NativeEvent != "beforeShellExecution" || e.Command != ".cursor/hooks/block-force-push.sh" || e.TimeoutSec != 20 {
		t.Fatalf("entry after round-trip: %+v", e)
	}
}

func TestParseClaude_extrasRoundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(claudeSettings), agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "MessageDisplay") || !strings.Contains(s, "plain.sh") {
		t.Errorf("MessageDisplay extra not restored: %s", s)
	}
	cfg2, _, err := Parse(out, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	assertClaudeConfigEqual(t, cfg1, cfg2)
}

func TestClaude_handlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "command", "command": "lint.sh", "async": true, "timeout": 15}]
	    }]
	  }
	}`
	cfg1, _, err := Parse([]byte(raw), agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[agenthooks.KindPreTool][0].Raw, "async", true) {
		t.Fatalf("async not preserved in Raw: %s", cfg1.Hooks[agenthooks.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"async": true`) {
		t.Fatalf("async not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	assertClaudeConfigEqual(t, cfg1, cfg2)
}

func TestClaude_claudeOnlyHandlerPreservesMatcher(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "mcp_tool", "tool": "my-server", "name": "list_items"}]
	    }]
	  }
	}`
	cfg1, warns, err := Parse([]byte(raw), agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg1.Extras) != 1 {
		t.Fatalf("want 1 extra, got %d hooks=%d warns=%v", len(cfg1.Extras), len(cfg1.Hooks), warns)
	}
	var g struct {
		Matcher string            `json:"matcher"`
		Hooks   []json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal(cfg1.Extras[0].Raw, &g); err != nil {
		t.Fatal(err)
	}
	if g.Matcher != "Bash" || len(g.Hooks) != 1 {
		t.Fatalf("extra group: matcher=%q hooks=%d", g.Matcher, len(g.Hooks))
	}
	out, _, err := Emit(cfg1, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"matcher": "Bash"`) || !strings.Contains(string(out), `"mcp_tool"`) {
		t.Fatalf("matcher or handler type lost: %s", out)
	}
	cfg2, _, err := Parse(out, agenthooks.Claude)
	if err != nil {
		t.Fatal(err)
	}
	assertClaudeConfigEqual(t, cfg1, cfg2)
}

func TestCopilot_handlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"command","command":"x.sh","cwd":"/tmp","timeoutSec":10}]}}`
	cfg1, _, err := Parse([]byte(raw), agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[agenthooks.KindPreTool][0].Raw, "cwd", "/tmp") {
		t.Fatalf("cwd not preserved: %s", cfg1.Hooks[agenthooks.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1, agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"cwd": "/tmp"`) {
		t.Fatalf("cwd not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out, agenthooks.Copilot)
	if err != nil {
		t.Fatal(err)
	}
	assertCopilotHooksEqual(t, cfg1, cfg2)
}

func TestCursor_handlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"command":"x.sh","loop_limit":3}]}}`
	cfg1, _, err := Parse([]byte(raw), agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[agenthooks.KindPreTool][0].Raw, "loop_limit", float64(3)) {
		t.Fatalf("loop_limit not preserved: %s", cfg1.Hooks[agenthooks.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"loop_limit": 3`) {
		t.Fatalf("loop_limit not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out, agenthooks.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	assertCursorHooksEqual(t, cfg1, cfg2)
}

func TestParse_unknownDialect(t *testing.T) {
	_, _, err := Parse([]byte("{}"), agenthooks.Unknown)
	if err == nil {
		t.Fatal("expected error for unknown dialect")
	}
}

func assertClaudeConfigEqual(t *testing.T, a, b Config) {
	t.Helper()
	if len(a.Extras) != len(b.Extras) {
		t.Fatalf("Extras len %d != %d", len(a.Extras), len(b.Extras))
	}
	for i := range a.Extras {
		if a.Extras[i].Event != b.Extras[i].Event {
			t.Fatalf("Extras[%d] event %q != %q", i, a.Extras[i].Event, b.Extras[i].Event)
		}
		if !rawJSONEqual(a.Extras[i].Raw, b.Extras[i].Raw) {
			t.Fatalf("Extras[%d] raw mismatch:\n%s\n!=\n%s", i, a.Extras[i].Raw, b.Extras[i].Raw)
		}
	}
	assertEntriesByKind(t, a.Hooks, b.Hooks)
}

func assertCopilotHooksEqual(t *testing.T, a, b Config) {
	t.Helper()
	assertEntriesByKind(t, a.Hooks, b.Hooks)
}

func assertCursorHooksEqual(t *testing.T, a, b Config) {
	t.Helper()
	assertEntriesByKind(t, a.Hooks, b.Hooks)
}

func assertEntriesByKind(t *testing.T, a, b map[agenthooks.Kind][]Entry) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("kind count %d != %d", len(a), len(b))
	}
	for kind, aEntries := range a {
		bEntries, ok := b[kind]
		if !ok || len(aEntries) != len(bEntries) {
			t.Fatalf("kind %q: len %d != %d", kind, len(aEntries), len(bEntries))
		}
		for i := range aEntries {
			ae, be := aEntries[i], bEntries[i]
			if ae.Command != be.Command || ae.Matcher != be.Matcher || ae.TimeoutSec != be.TimeoutSec ||
				ae.NativeEvent != be.NativeEvent || ae.Type != be.Type {
				t.Fatalf("kind %q entry %d: %+v != %+v", kind, i, ae, be)
			}
			if !rawJSONEqual(ae.ClaudeGroupIf, be.ClaudeGroupIf) {
				t.Fatalf("kind %q entry %d ClaudeGroupIf mismatch:\n%s\n!=\n%s", kind, i, ae.ClaudeGroupIf, be.ClaudeGroupIf)
			}
			if !rawJSONEqual(ae.Raw, be.Raw) {
				t.Fatalf("kind %q entry %d raw mismatch:\n%s\n!=\n%s", kind, i, ae.Raw, be.Raw)
			}
		}
	}
}

func rawJSONEqual(a, b json.RawMessage) bool {
	if bytes.Equal(a, b) {
		return true
	}
	var va, vb any
	if json.Unmarshal(a, &va) != nil || json.Unmarshal(b, &vb) != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

func jsonFieldEqual(t *testing.T, raw json.RawMessage, key string, want any) bool {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	got, ok := m[key]
	if !ok {
		return false
	}
	return reflect.DeepEqual(got, want)
}

func TestEmit_producesValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		dialect agenthooks.Dialect
	}{
		{name: "claude", raw: claudeSettings, dialect: agenthooks.Claude},
		{name: "copilot", raw: copilotSettings, dialect: agenthooks.Copilot},
		{name: "cursor", raw: cursorSettings, dialect: agenthooks.Cursor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _, err := Parse([]byte(tt.raw), tt.dialect)
			if err != nil {
				t.Fatal(err)
			}
			out, _, err := Emit(cfg, tt.dialect)
			if err != nil {
				t.Fatal(err)
			}
			if !json.Valid(out) {
				t.Fatalf("invalid JSON: %s", out)
			}
		})
	}
}
