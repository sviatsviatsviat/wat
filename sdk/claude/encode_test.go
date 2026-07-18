package claude_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
)

func TestEncode_PreToolDeny(t *testing.T) {
	out, err := claude.Encode("PreToolUse", claude.PreToolUseResultsForTest().Deny("destructive command"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_UserPromptBlock(t *testing.T) {
	out, err := claude.Encode("UserPromptSubmit", claude.UserPromptSubmitResultsForTest().Block("blocked prompt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopBlock(t *testing.T) {
	out, err := claude.Encode("Stop", claude.StopResultsForTest().FollowUp("run the tests"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_SessionStartEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.sh")
	var written []byte
	out, err := claude.Encode("SessionStart", claude.SessionStartResultsForTest().Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}),
		claude.WithGetenv(func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return envPath
			}
			return ""
		}),
		claude.WithAppendFile(func(path string, data []byte) error {
			written = append(written, data...)
			return os.WriteFile(path, written, 0o644)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("env-only result should produce no stdout, got %q", out)
	}
	got := string(written)
	for _, line := range []string{`export FOO='bar'`, `export BAZ='qux'`} {
		if !strings.Contains(got, line) {
			t.Fatalf("env file = %q, want lines containing %s", got, line)
		}
	}
}

func TestEncode_ZeroOutput(t *testing.T) {
	out, err := claude.Encode("PreToolUse", nil)
	if err != nil || out != nil {
		t.Fatalf("zero output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, err := claude.Encode(claude.EventPostToolUse, claude.PreToolUseResultsForTest().Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestEncode_PointerOutput(t *testing.T) {
	deny := claude.PreToolUseResultsForTest().Deny("blocked")
	out, err := claude.Encode("PreToolUse", &deny)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}

	var typedNil *claude.PreToolUseOutput
	out, err = claude.Encode("PreToolUse", typedNil)
	if err != nil || out != nil {
		t.Fatalf("nil pointer output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, err := claude.Encode("PermissionRequest", claude.PermissionRequestResultsForTest().Deny("policy").WithInterrupt(true))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), `"continue":false`) {
		t.Fatalf("interrupt must not set top-level continue: %s", out)
	}
	if !strings.Contains(string(out), `"interrupt":true`) {
		t.Fatalf("bad output: %s", out)
	}
}
