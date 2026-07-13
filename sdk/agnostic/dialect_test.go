package agnostic

import "testing"

const claudePreToolUse = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "PreToolUse",
  "tool_name": "Bash",
  "tool_use_id": "tu_1",
  "tool_input": {"command": "rm -rf /tmp/build", "description": "clean"}
}`

const copilotCamelPreToolUse = `{
  "sessionId": "s1",
  "timestamp": 1760000000000,
  "cwd": "/w",
  "toolName": "bash",
  "toolArgs": {"command": "rm -rf /"}
}`

const copilotVSCodeStop = `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "transcript_path": "/tmp/t",
  "stop_reason": "end_turn"
}`

// Minimal Cursor payloads for dialect detection only (no cursor codec dependency).
const detectCursorShell = `{"conversation_id":"c1","cursor_version":"1.7.2","hook_event_name":"beforeShellExecution"}`
const detectCursorStop = `{"conversation_id":"c1","cursor_version":"1.7.2","hook_event_name":"stop"}`

func TestParseDialect(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Dialect
	}{
		{name: "claude", in: "claude", want: Claude},
		{name: "claude-code alias", in: "claude-code", want: Claude},
		{name: "claudecode alias", in: "claudecode", want: Claude},
		{name: "copilot", in: "copilot", want: Copilot},
		{name: "github-copilot alias", in: "github-copilot", want: Copilot},
		{name: "gh alias", in: "gh", want: Copilot},
		{name: "cursor", in: "cursor", want: Cursor},
		{name: "trimmed whitespace", in: "  CURSOR  ", want: Cursor},
		{name: "unknown", in: "vscode", want: Unknown},
		{name: "empty", in: "", want: Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDialect(tt.in); got != tt.want {
				t.Fatalf("ParseDialect(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestDetect_payloadFixtures(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    Dialect
	}{
		{name: "claude pre tool use", payload: claudePreToolUse, want: Claude},
		{name: "copilot camelCase", payload: copilotCamelPreToolUse, want: Copilot},
		{name: "copilot VS Code format", payload: copilotVSCodeStop, want: Copilot},
		{name: "cursor shell", payload: detectCursorShell, want: Cursor},
		{name: "cursor stop", payload: detectCursorStop, want: Cursor},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Detect([]byte(tt.payload), nil); got != tt.want {
				t.Fatalf("Detect = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetect_ambiguousEnv(t *testing.T) {
	t.Run("cursor exports CLAUDE_PROJECT_DIR", func(t *testing.T) {
		env := map[string]string{
			"CLAUDE_PROJECT_DIR": "/p",
			"CURSOR_VERSION":     "1.7",
		}
		got := Detect([]byte(`{}`), func(k string) string { return env[k] })
		if got != Cursor {
			t.Fatalf("Detect = %v, want Cursor", got)
		}
	})

	t.Run("copilot payload wins over CLAUDE_PROJECT_DIR", func(t *testing.T) {
		env := map[string]string{"CLAUDE_PROJECT_DIR": "/p"}
		got := Detect([]byte(copilotVSCodeStop), func(k string) string { return env[k] })
		if got != Copilot {
			t.Fatalf("Detect = %v, want Copilot", got)
		}
	})

	t.Run("env-only CLAUDE_PROJECT_DIR", func(t *testing.T) {
		env := map[string]string{"CLAUDE_PROJECT_DIR": "/p"}
		got := Detect([]byte(`{}`), func(k string) string { return env[k] })
		if got != Claude {
			t.Fatalf("Detect = %v, want Claude", got)
		}
	})

	t.Run("COPILOT_HOME env fallback", func(t *testing.T) {
		env := map[string]string{"COPILOT_HOME": "/home/user/.copilot"}
		got := Detect([]byte(`{}`), func(k string) string { return env[k] })
		if got != Copilot {
			t.Fatalf("Detect = %v, want Copilot", got)
		}
	})
}

func TestDetect_edgeCases(t *testing.T) {
	t.Run("nil getenv", func(t *testing.T) {
		if got := Detect([]byte(claudePreToolUse), nil); got != Claude {
			t.Fatalf("Detect = %v, want Claude", got)
		}
	})

	t.Run("invalid JSON env fallback", func(t *testing.T) {
		env := map[string]string{"CURSOR_VERSION": "1.7"}
		got := Detect([]byte("not json"), func(k string) string { return env[k] })
		if got != Cursor {
			t.Fatalf("Detect = %v, want Cursor", got)
		}
	})

	t.Run("empty payload no env", func(t *testing.T) {
		if got := Detect([]byte(`{}`), nil); got != Unknown {
			t.Fatalf("Detect = %v, want Unknown", got)
		}
	})
}
