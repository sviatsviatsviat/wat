package portconfig

import (
	"encoding/json"
	"strings"
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func TestTranslate_claudeToCursor(t *testing.T) {
	out, warns, err := Translate([]byte(claudeSettings), sdkclaude.Dialect, sdkcursor.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Version int `json:"version"`
		Hooks   map[string][]struct {
			Command string `json:"command"`
			Matcher string `json:"matcher"`
			Timeout int    `json:"timeout"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(out, &f); err != nil {
		t.Fatal(err)
	}
	if f.Version != 1 {
		t.Fatal("cursor config must be version 1")
	}
	pre := f.Hooks["preToolUse"]
	if len(pre) != 2 {
		t.Fatalf("want 2 preToolUse handlers, got %d: %s", len(pre), out)
	}
	var matchers []string
	for _, h := range pre {
		matchers = append(matchers, h.Matcher)
	}
	joined := strings.Join(matchers, " ")
	if !strings.Contains(joined, "Shell") {
		t.Errorf("Bash should map to Shell, got %q", joined)
	}
	if len(f.Hooks["stop"]) != 1 {
		t.Errorf("Stop should map to stop: %s", out)
	}
	if _, ok := f.Hooks["MessageDisplay"]; ok {
		t.Error("MessageDisplay must not survive into Cursor config")
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "MessageDisplay") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a warning about MessageDisplay, got %v", warns)
	}
}

func TestTranslate_claudeToCopilot(t *testing.T) {
	out, _, err := Translate([]byte(claudeSettings), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Version int `json:"version"`
		Hooks   map[string][]struct {
			Type       string `json:"type"`
			Command    string `json:"command"`
			Matcher    string `json:"matcher"`
			TimeoutSec int    `json:"timeoutSec"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(out, &f); err != nil {
		t.Fatal(err)
	}
	pre := f.Hooks["PreToolUse"]
	if len(pre) != 2 {
		t.Fatalf("want 2 PreToolUse handlers: %s", out)
	}
	foundBash := false
	for _, h := range pre {
		if strings.Contains(h.Matcher, "bash") {
			foundBash = true
		}
		if strings.Contains(h.Matcher, "bash") && h.TimeoutSec != 15 {
			t.Errorf("timeout should carry over as timeoutSec, got %d", h.TimeoutSec)
		}
	}
	if !foundBash {
		t.Errorf("Bash should map to bash: %s", out)
	}
	if len(f.Hooks["Stop"]) != 1 {
		t.Errorf("Stop should map to Stop: %s", out)
	}
}

func TestTranslate_cursorToClaude(t *testing.T) {
	out, _, err := Translate([]byte(cursorSettings), sdkcursor.Dialect, sdkclaude.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "UserPromptSubmit") {
		t.Errorf("beforeSubmitPrompt should map to UserPromptSubmit: %s", s)
	}
	if !strings.Contains(s, `"matcher": "Bash"`) {
		t.Errorf("Shell should map back to Bash: %s", s)
	}
}

func TestTranslate_copilotAnchoredRegex(t *testing.T) {
	t.Run("copilot_to_claude", func(t *testing.T) {
		raw := `{"version":1,"hooks":{"PreToolUse":[{"type":"command","command":"x.sh","matcher":"^(?:bash)$"}]}}`
		out, warns, err := Translate([]byte(raw), sdkcopilot.Dialect, sdkclaude.Dialect)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `"matcher": "Bash"`) {
			t.Fatalf("anchored regex should un-anchor and map bash to Bash for Claude: %s", out)
		}
		found := false
		for _, w := range warns {
			if strings.Contains(string(w), "un-anchored") {
				found = true
			}
		}
		if !found {
			t.Errorf("expected un-anchor warning, got %v", warns)
		}
	})
	t.Run("claude_to_copilot", func(t *testing.T) {
		raw := `{"hooks":{"PreToolUse":[{"matcher":"bash|edit","hooks":[{"type":"command","command":"x.sh"}]}]}}`
		out, _, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(out), `"matcher": "^(?:bash|edit)$"`) {
			t.Fatalf("simple alternation should anchor for Copilot: %s", out)
		}
	})
}

func TestTranslate_claudeIfRule(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "command", "command": "guard.sh", "if": {"permissionMode": "default"}}]
	    }]
	  }
	}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcursor.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "guard.sh") {
		t.Fatalf("command should still port: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "if permission rule") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected if permission rule warning, got %v", warns)
	}
}

func TestTranslate_unsupportedEvent(t *testing.T) {
	raw := `{
	  "hooks": {
	    "Notification": [
	      {"hooks": [{"type": "command", "command": "notify.sh"}]}
	    ]
	  }
	}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcursor.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "Notification") {
		t.Errorf("Notification must not appear in Cursor config: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "Notification") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Notification warning, got %v", warns)
	}
}

func TestTranslate_sameDialect(t *testing.T) {
	out, warns, err := Translate([]byte(claudeSettings), sdkclaude.Dialect, sdkclaude.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != claudeSettings {
		t.Fatal("same-dialect translate should return input unchanged")
	}
	if len(warns) != 0 {
		t.Fatalf("same-dialect translate should produce no warnings, got %v", warns)
	}
}

func TestTranslate_httpToCursor(t *testing.T) {
	raw := `{"hooks":{"PreToolUse":[{"matcher":"Bash","hooks":[{"type":"http","url":"https://example.com/hook"}]}]}}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcursor.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "example.com") {
		t.Errorf("http handler should not appear in Cursor config: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "http") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected http warning, got %v", warns)
	}
}

func TestTranslate_timeoutDefaultWarning(t *testing.T) {
	raw := `{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"done.sh"}]}]}}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "600") && strings.Contains(string(w), "30") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected default timeout warning, got %v", warns)
	}
	if !strings.Contains(string(out), `"timeoutSec": 600`) {
		t.Errorf("expected explicit source default timeout in output: %s", out)
	}
}

func TestTranslate_mcpToolExtraNotEmitted(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "hooks": [{"type": "mcp_tool", "tool": "my-server", "name": "list_items"}]
	    }]
	  }
	}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "mcp_tool") {
		t.Fatalf("mcp_tool extra must not appear in Copilot config: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "PreToolUse") && strings.Contains(string(w), "not ported") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected PreToolUse extra warning, got %v", warns)
	}
}

func TestTranslate_cursorDedicatedEvent(t *testing.T) {
	out, warns, err := Translate([]byte(cursorDedicatedSettings), sdkcursor.Dialect, sdkclaude.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "PreToolUse") {
		t.Fatalf("beforeShellExecution should map to PreToolUse: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "beforeShellExecution") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected dedicated event warning, got %v", warns)
	}
}

func TestTranslate_droppedRawFieldsWarning(t *testing.T) {
	raw := `{"hooks":{"PreToolUse":[{"matcher":"Bash","hooks":[{"type":"command","command":"x.sh","async":true}]}]}}`
	_, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "async") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected dropped field warning for async, got %v", warns)
	}
}

func TestTranslate_claudeGroupIfRule(t *testing.T) {
	raw := `{
	  "hooks": {
	    "PreToolUse": [{
	      "matcher": "Bash",
	      "if": {"permissionMode": "default"},
	      "hooks": [{"type": "command", "command": "guard.sh"}]
	    }]
	  }
	}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "guard.sh") {
		t.Fatalf("command should still port: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "if permission rule") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected group if permission rule warning, got %v", warns)
	}
}

func TestTranslate_complexRegexKeptVerbatim(t *testing.T) {
	raw := `{"hooks":{"PreToolUse":[{"matcher":"Bash{1,3}","hooks":[{"type":"command","command":"x.sh"}]}]}}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "Bash{1,3}") {
		t.Fatalf("complex regex should stay unchanged: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "complex regex kept verbatim") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected complex regex warning, got %v", warns)
	}
}

func TestTranslate_sourceNamedRegexWarns(t *testing.T) {
	raw := `{"hooks":{"PreToolUse":[{"matcher":"^bash$","hooks":[{"type":"command","command":"x.sh"}]}]}}`
	out, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "^bash$") {
		t.Fatalf("source-named regex should stay unchanged: %s", out)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), "complex regex kept verbatim") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected complex regex warning for ^bash$, got %v", warns)
	}
}

func TestTranslate_unknownTokenAlwaysWarns(t *testing.T) {
	raw := `{"hooks":{"PreToolUse":[{"matcher":"unknownTool","hooks":[{"type":"command","command":"x.sh"}]}]}}`
	_, warns, err := Translate([]byte(raw), sdkclaude.Dialect, sdkcopilot.Dialect)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, w := range warns {
		if strings.Contains(string(w), `matcher token "unknownTool"`) {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning for unknown token, got %v", warns)
	}
}
