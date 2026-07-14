package cursor

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
)

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

func TestParse_dedicatedEvent(t *testing.T) {
	cfg, _, err := Parse([]byte(cursorDedicatedSettings))
	if err != nil {
		t.Fatal(err)
	}
	pre := cfg.Hooks[agnostic.KindPreTool]
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

func TestEmit_roundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(cursorSettings))
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertHooksEqual(t, cfg1, cfg2)
}

func TestEmit_dedicatedEventRoundTrip(t *testing.T) {
	cfg1, _, err := Parse([]byte(cursorDedicatedSettings))
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `"beforeShellExecution"`) {
		t.Errorf("beforeShellExecution not preserved: %s", s)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg2.Hooks[agnostic.KindPreTool]) != 1 {
		t.Fatalf("want 1 pre-tool entry after round-trip, got %d", len(cfg2.Hooks[agnostic.KindPreTool]))
	}
	e := cfg2.Hooks[agnostic.KindPreTool][0]
	if e.NativeEvent != "beforeShellExecution" || e.Command != ".cursor/hooks/block-force-push.sh" || e.TimeoutSec != 20 {
		t.Fatalf("entry after round-trip: %+v", e)
	}
}

func TestHandlerExtraFieldsRoundTrip(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"command":"x.sh","loop_limit":3}]}}`
	cfg1, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !jsonFieldEqual(t, cfg1.Hooks[agnostic.KindPreTool][0].Raw, "loop_limit", float64(3)) {
		t.Fatalf("loop_limit not preserved: %s", cfg1.Hooks[agnostic.KindPreTool][0].Raw)
	}
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"loop_limit": 3`) {
		t.Fatalf("loop_limit not in emitted JSON: %s", out)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	assertHooksEqual(t, cfg1, cfg2)
}

func TestEmit_mutateRawOverlay(t *testing.T) {
	raw := `{"version":1,"hooks":{"preToolUse":[{"command":"x.sh","matcher":"Shell","loop_limit":3,"timeout":10}]}}`
	cfg1, _, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	entry := &cfg1.Hooks[agnostic.KindPreTool][0]
	entry.TimeoutSec = 42
	entry.Matcher = "Read"
	out, _, err := Emit(cfg1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"loop_limit": 3`) {
		t.Fatalf("loop_limit not preserved: %s", out)
	}
	if !strings.Contains(string(out), `"timeout": 42`) {
		t.Fatalf("mutated timeout not emitted: %s", out)
	}
	if !strings.Contains(string(out), `"matcher": "Read"`) {
		t.Fatalf("mutated matcher not emitted: %s", out)
	}
	cfg2, _, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if cfg2.Hooks[agnostic.KindPreTool][0].TimeoutSec != 42 {
		t.Fatalf("timeout not round-tripped: %+v", cfg2.Hooks[agnostic.KindPreTool][0])
	}
	if !jsonFieldEqual(t, cfg2.Hooks[agnostic.KindPreTool][0].Raw, "loop_limit", float64(3)) {
		t.Fatalf("loop_limit lost from raw: %s", cfg2.Hooks[agnostic.KindPreTool][0].Raw)
	}
}

func TestEmit_rejectsUnsupportedHandlerType(t *testing.T) {
	cfg := model.Config{
		Hooks: map[agnostic.Kind][]model.Entry{
			agnostic.KindPreTool: {{
				Kind:        agnostic.KindPreTool,
				NativeEvent: "preToolUse",
				Type:        "http",
				URL:         "https://example.com/hook",
			}},
		},
	}
	out, warns, err := Emit(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "example.com") {
		t.Fatalf("http handler should be dropped: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "http") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected http drop warning, got %v", warns)
	}
}

func assertHooksEqual(t *testing.T, a, b model.Config) {
	t.Helper()
	assertEntriesByKind(t, a.Hooks, b.Hooks)
}

func assertEntriesByKind(t *testing.T, a, b map[agnostic.Kind][]model.Entry) {
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
