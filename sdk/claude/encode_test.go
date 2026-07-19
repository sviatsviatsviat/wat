package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncode_PreToolDeny(t *testing.T) {
	out, err := encode("PreToolUse", preToolUseResults{}.Deny("destructive command"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_UserPromptBlock(t *testing.T) {
	out, err := encode("UserPromptSubmit", userPromptSubmitResults{}.Block("blocked prompt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopBlock(t *testing.T) {
	out, err := encode("Stop", stopResults{}.FollowUp("run the tests"))
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
	out, err := encode("SessionStart", sessionStartResults{}.Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}),
		WithGetenv(func(key string) string {
			if key == "CLAUDE_ENV_FILE" {
				return envPath
			}
			return ""
		}),
		WithAppendFile(func(path string, data []byte) error {
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
	out, err := encode("PreToolUse", nil)
	if err != nil || out != nil {
		t.Fatalf("zero output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, err := encode(EventPostToolUse, preToolUseResults{}.Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestEncode_NilOutput(t *testing.T) {
	var typedNil PreToolUseOutput
	out, err := encode("PreToolUse", typedNil)
	if err != nil || out != nil {
		t.Fatalf("nil output should be silent, got %q err=%v", out, err)
	}
}

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, err := encode("PermissionRequest", permissionRequestResults{}.Deny("policy").WithInterrupt(true))
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

func TestWriteEnvFile_InvalidKey(t *testing.T) {
	err := WriteEnvFile(
		map[string]string{"FOO\nBAR": "value"},
		func(string) string { return "/tmp/env.sh" },
		nil,
	)
	if err == nil {
		t.Fatal("expected error for invalid env key")
	}
}
