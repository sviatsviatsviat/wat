package portconfig

import (
	"encoding/json"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
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

func TestParse_unknownDialect(t *testing.T) {
	_, _, err := Parse([]byte("{}"), agnostic.Unknown)
	if err == nil {
		t.Fatal("expected error for unknown dialect")
	}
}

func TestEmit_producesValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		dialect agnostic.Dialect
	}{
		{name: "claude", raw: claudeSettings, dialect: agnostic.Claude},
		{name: "copilot", raw: copilotSettings, dialect: agnostic.Copilot},
		{name: "cursor", raw: cursorSettings, dialect: agnostic.Cursor},
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
