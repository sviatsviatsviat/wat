package claude

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncode_PreToolDeny(t *testing.T) {
	out, code, err := codec.Encode("PreToolUse", preToolUseResults{}.Deny("destructive command"))
	if err != nil {
		t.Fatal(err)
	}
	if code != SuccessExit {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(string(out), `"permissionDecision":"deny"`) {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_UserPromptBlock(t *testing.T) {
	out, _, err := codec.Encode("UserPromptSubmit", userPromptSubmitResults{}.Block("blocked prompt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "blocked prompt") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_StopBlock(t *testing.T) {
	out, _, err := codec.Encode("Stop", stopResults{}.FollowUp("run the tests"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"decision":"block"`) || !strings.Contains(string(out), "run the tests") {
		t.Fatalf("bad output: %s", out)
	}
}

func TestEncode_SessionStartEnv(t *testing.T) {
	envPath := filepath.Join(t.TempDir(), "env.sh")
	t.Setenv("CLAUDE_ENV_FILE", envPath)
	out, code, err := codec.Encode("SessionStart", sessionStartResults{}.Noop().WithEnv(map[string]string{"FOO": "bar", "BAZ": "qux"}))
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

func TestEncode_ZeroOutput(t *testing.T) {
	out, code, err := codec.Encode("PreToolUse", preToolUseResults{}.Noop())
	if err != nil || out != nil || code != SuccessExit {
		t.Fatalf("zero output should be silent, got %q code=%d err=%v", out, code, err)
	}
}

func TestEncode_EventOutputMismatch(t *testing.T) {
	_, _, err := codec.Encode(EventPostToolUse, preToolUseResults{}.Allow())
	if err == nil {
		t.Fatal("expected incompatible event/output error")
	}
}

func TestEncode_PermissionRequestInterrupt(t *testing.T) {
	out, _, err := codec.Encode("PermissionRequest", permissionRequestResults{}.Deny("policy").WithInterrupt(true))
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
	err := writeEnvFile(
		map[string]string{"FOO\nBAR": "value"},
		func(string) string { return "/tmp/env.sh" },
		nil,
	)
	if err == nil {
		t.Fatal("expected error for invalid env key")
	}
}
