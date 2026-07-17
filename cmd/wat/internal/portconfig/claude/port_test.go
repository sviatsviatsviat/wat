package claude

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
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

// TestParse_mappableAndExtras verifies Claude settings map known events and preserve unmappable extras.
func TestParse_mappableAndExtras(t *testing.T) {
	cfg, _, err := Parse([]byte(claudeSettings))
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[model.KindPreTool]
	if len(pre) != 2 {
		t.Fatalf("want 2 PreToolUse entries, got %d", len(pre))
	}
	if pre[0].Command != ".claude/hooks/block-rm.sh" || pre[0].Matcher != "Bash" || pre[0].TimeoutSec != 15 {
		t.Fatalf("first PreToolUse entry: %+v", pre[0])
	}
	if len(cfg.Hooks[model.KindStop]) != 1 {
		t.Fatalf("want 1 Stop entry, got %d", len(cfg.Hooks[model.KindStop]))
	}
	if len(cfg.Extras) != 1 || cfg.Extras[0].Event != "MessageDisplay" {
		t.Fatalf("MessageDisplay should be in Extras: %+v", cfg.Extras)
	}
}

// TestEmit_roundTrip checks that mutating a parsed entry survives emit and re-parse.
func TestEmit_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(claudeSettings))
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg1.Hooks[model.KindPreTool]
	if len(pre) == 0 {
		t.Fatal("expected PreTool entries")
	}
	wantCommand := pre[0].Command
	pre[0].TimeoutSec = 99
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Hooks[model.KindPreTool][0].TimeoutSec != 99 {
		t.Fatalf("mutated timeout not emitted: %+v", cfg2.Hooks[model.KindPreTool][0])
	}
	emittedRaw := cfg2.Hooks[model.KindPreTool][0].Raw
	if !jsonFieldEqual(t, emittedRaw, "timeout", float64(99)) {
		t.Fatalf("timeout not overlaid in raw JSON: %s", emittedRaw)
	}
	if !jsonFieldEqual(t, emittedRaw, "command", wantCommand) {
		t.Fatalf("command not preserved in raw JSON: %s", emittedRaw)
	}
	cfg1.Hooks[model.KindPreTool][0].TimeoutSec = cfg2.Hooks[model.KindPreTool][0].TimeoutSec
	cfg1.Hooks[model.KindPreTool][0].Raw = nil
	cfg2.Hooks[model.KindPreTool][0].Raw = nil
	assertConfigEqual(t, cfg1, cfg2)
}

// TestParse_extrasRoundTrip verifies unmappable Claude events round-trip through emit.
func TestParse_extrasRoundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(claudeSettings))
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "MessageDisplay") || !strings.Contains(s, "plain.sh") {
		t.Errorf("MessageDisplay extra not restored: %s", s)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertConfigEqual(t, cfg1, cfg2)
}

// TestHandlerExtraFieldsRoundTrip verifies unknown handler JSON fields survive Claude emit.
func TestHandlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "command", "command": "lint.sh", "async": true, "timeout": 15}]
	    }]
	  }
	}`
	cfg1, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[model.KindPreTool][0].Raw, "async", true) {
		t.Fatalf("async not preserved in Raw: %s", cfg1.Hooks[model.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"async": true`) {
		t.Fatalf("async not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertConfigEqual(t, cfg1, cfg2)
}

// TestClaudeOnlyHandlerPreservesMatcher verifies Claude-only handler types keep matcher groups in Extras.
func TestClaudeOnlyHandlerPreservesMatcher(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "mcp_tool", "tool": "my-server", "name": "list_items"}]
	    }]
	  }
	}`
	cfg1, warns, err := Parse([]byte(raw))
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
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"matcher": "Bash"`) || !strings.Contains(string(out), `"mcp_tool"`) {
		t.Fatalf("matcher or handler type lost: %s", out)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertConfigEqual(t, cfg1, cfg2)
}

func assertConfigEqual(t *testing.T, a, b model.Config) {
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

func assertEntriesByKind(t *testing.T, a, b map[model.Kind][]model.Entry) {
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
