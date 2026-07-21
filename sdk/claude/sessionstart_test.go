package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncode_SessionStartEnv(t *testing.T) {
	envPath := filepath.Join(t.TempDir(), "env.sh")
	t.Setenv("CLAUDE_ENV_FILE", envPath)
	out, code, err := sessionStartResults{}.Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
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

func TestDecode_SessionStart(t *testing.T) {
	mustDecode[SessionStart](t, `{"session_id":"s","hook_event_name":"SessionStart","source":"startup","model":"claude-3"}`, EventSessionStart)
}
