package copilot

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
)

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

func TestParse_invalidHandlerPreservedInExtras(t *testing.T) {
	const copilotInvalid = `{
  "version": 1,
  "hooks": {
    "preToolUse": [{"command": ["not-a-string"]}]
  }
}`
	cfg, warns, err := Parse([]byte(copilotInvalid))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Hooks[model.KindPreTool]) != 0 {
		t.Fatalf("invalid handler should not map: %+v", cfg.Hooks[model.KindPreTool])
	}
	if len(cfg.Extras) != 1 {
		t.Fatalf("want 1 extra, got %d", len(cfg.Extras))
	}
	if cfg.Extras[0].Event != "preToolUse" {
		t.Fatalf("extra event = %q", cfg.Extras[0].Event)
	}
	wantRaw := json.RawMessage(`{"command": ["not-a-string"]}`)
	if !rawJSONEqual(cfg.Extras[0].Raw, wantRaw) {
		t.Fatalf("extra raw = %s, want %s", cfg.Extras[0].Raw, wantRaw)
	}
	if len(warns) == 0 || !strings.Contains(string(warns[0]), "invalid handler JSON") {
		t.Fatalf("want parse warning, got %v", warns)
	}

	out, emitWarns, err := Emit(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(emitWarns) != 0 {
		t.Fatalf("emit warnings = %v", emitWarns)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg2.Extras) != 1 || !rawJSONEqual(cfg2.Extras[0].Raw, wantRaw) {
		t.Fatalf("round-trip extra = %+v", cfg2.Extras)
	}
}

func TestParse_timeoutAndCommand(t *testing.T) {
	cfg, _, err := Parse([]byte(copilotSettings))
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[model.KindPreTool]
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

func TestParse_timeoutAlias(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"command","command":"x.sh","timeout":42}]}}`
	cfg, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Hooks[model.KindPreTool][0].TimeoutSec != 42 {
		t.Fatalf("timeout alias not resolved: %+v", cfg.Hooks[model.KindPreTool][0])
	}
}

func TestParse_powershellOnlyNotCommand(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"command","powershell":"Get-ChildItem"}]}}`
	cfg, warns, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if len(warns) == 0 {
		t.Fatal("expected powershell warning")
	}
	powershellWarn := false
	for _, w := range warns {
		if strings.Contains(string(w), "powershell") {
			powershellWarn = true
			break
		}
	}
	if !powershellWarn {
		t.Fatalf("expected warning mentioning powershell, got %v", warns)
	}
	entry := cfg.Hooks[model.KindPreTool][0]
	if entry.Command != "" {
		t.Fatalf("Command = %q, want empty", entry.Command)
	}
	out, _, err := Emit(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"powershell"`) || !strings.Contains(string(out), "Get-ChildItem") {
		t.Fatalf("emit lost powershell field: %s", out)
	}
}

func TestEmit_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(copilotSettings))
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg1.Hooks[model.KindPreTool]
	if len(pre) == 0 {
		t.Fatal("expected PreTool entries")
	}
	for i := range pre {
		if pre[i].Matcher == "bash" {
			pre[i].TimeoutSec = 99
			break
		}
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range cfg2.Hooks[model.KindPreTool] {
		if e.Matcher == "bash" && e.TimeoutSec != 99 {
			t.Fatalf("mutated timeout not emitted: %+v", e)
		}
	}
	for i := range cfg1.Hooks[model.KindPreTool] {
		for j := range cfg2.Hooks[model.KindPreTool] {
			if cfg1.Hooks[model.KindPreTool][i].Matcher == cfg2.Hooks[model.KindPreTool][j].Matcher {
				cfg1.Hooks[model.KindPreTool][i].Raw = cfg2.Hooks[model.KindPreTool][j].Raw
			}
		}
	}
	assertHooksEqual(t, cfg1, cfg2)
}

func TestHandlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"command","command":"x.sh","cwd":"/tmp","timeoutSec":10}]}}`
	cfg1, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[model.KindPreTool][0].Raw, "cwd", "/tmp") {
		t.Fatalf("cwd not preserved: %s", cfg1.Hooks[model.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"cwd": "/tmp"`) {
		t.Fatalf("cwd not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertHooksEqual(t, cfg1, cfg2)
}

func TestEmit_clearsStaleFieldsOnOverlay(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"type":"http","url":"https://old.example","matcher":"bash","timeoutSec":5,"cwd":"/tmp"}]}}`
	cfg1, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	entry := &cfg1.Hooks[model.KindPreTool][0]
	entry.Type = "command"
	entry.Command = "new.sh"
	entry.URL = ""
	entry.Matcher = ""
	entry.TimeoutSec = 0
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if strings.Contains(s, "old.example") || strings.Contains(s, `"url"`) {
		t.Fatalf("stale http url not removed: %s", s)
	}
	if !strings.Contains(s, `"cwd"`) {
		t.Fatalf("cwd should be preserved from raw: %s", s)
	}
	if !strings.Contains(s, `"command": "new.sh"`) {
		t.Fatalf("command overlay missing: %s", s)
	}
	if strings.Contains(s, `"matcher"`) {
		t.Fatalf("cleared matcher should be absent: %s", s)
	}
	if strings.Contains(s, `"timeoutSec"`) || strings.Contains(s, `"timeout"`) {
		t.Fatalf("cleared timeout should be absent: %s", s)
	}
}

func assertHooksEqual(t *testing.T, a, b model.Config) {
	t.Helper()
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
