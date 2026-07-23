package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func init() {
	Register(runtime.Codec)
}

func mustDecode[E any](t *testing.T, raw, wantName string) E {
	t.Helper()
	ev, err := runtime.Codec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != wantName {
		t.Fatalf("EventName() = %q, want %q", ev.EventName(), wantName)
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestEncode_Env(t *testing.T) {
	envPath := filepath.Join(t.TempDir(), "env.sh")
	t.Setenv("CLAUDE_ENV_FILE", envPath)
	out, code, err := results{}.Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != event.SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	if out != nil {
		t.Fatalf("env-only result should produce no stdout, got %q", out)
	}
	written, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(written)
	for _, line := range []string{`export FOO='bar'`, `export BAZ='qux'`} {
		if !strings.Contains(got, line) {
			t.Fatalf("env file = %q, want lines containing %s", got, line)
		}
	}
}

func TestDecode(t *testing.T) {
	ev := mustDecode[Event](t, `{"session_id":"s","hook_event_name":"SessionStart","source":"startup","model":"claude-3"}`, event.SessionStart)
	if ev.Source != "startup" || ev.Model != "claude-3" {
		t.Fatalf("session start fields = %+v", ev)
	}
}

func TestMerge_contextJoins(t *testing.T) {
	a := results{}.Context("one")
	b := results{}.Context("two")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if out.additionalContext != "one\n\ntwo" {
		t.Fatalf("context = %q", out.additionalContext)
	}
	if merged.Stop() {
		t.Fatal("context merge should not stop")
	}
}
