package checks

import (
	"testing"

	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
)

func TestParseGoModDirective(t *testing.T) {
	got, err := ParseGoModDirective([]byte("module wat-hooks\n\ngo 1.26\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.26" {
		t.Fatalf("got %q want 1.26", got)
	}
}

func TestGoVersionAtLeast(t *testing.T) {
	tests := []struct {
		installed string
		required  string
		want      bool
	}{
		{"1.26.0", "1.26", true},
		{"1.26.1", "1.26", true},
		{"1.25.9", "1.26", false},
		{"1.26", "1.26.0", true},
	}
	for _, tt := range tests {
		if got := GoVersionAtLeast(tt.installed, tt.required); got != tt.want {
			t.Fatalf("GoVersionAtLeast(%q, %q) = %v, want %v", tt.installed, tt.required, got, tt.want)
		}
	}
}

func TestParseWatRunFlags(t *testing.T) {
	agent, event, ok := ParseWatRunFlags("/usr/bin/wat run --agent claude --event PreToolUse")
	if !ok || agent != "claude" || event != "PreToolUse" {
		t.Fatalf("got %q %q %v", agent, event, ok)
	}
}

func TestIsWatManagedCommand_exactEvent(t *testing.T) {
	watAbs := "/usr/bin/wat"
	if !IsWatManagedCommand(watAbs+" run --agent claude --event PostToolUse", "claude", "PostToolUse", watAbs) {
		t.Fatal("expected match for PostToolUse")
	}
	if IsWatManagedCommand(watAbs+" run --agent claude --event PostToolUseFailure", "claude", "PostToolUse", watAbs) {
		t.Fatal("PostToolUseFailure must not match PostToolUse")
	}
	if IsWatManagedCommand(watAbs+" run --agent claude --event PostToolUse && echo", "claude", "PostToolUse", watAbs) {
		t.Fatal("trailing shell tokens must not match")
	}
}

func TestParseVersionParts_devSuffix(t *testing.T) {
	got := parseVersionParts("1.22-b9a08f159d")
	if len(got) != 2 || got[0] != 1 || got[1] != 22 {
		t.Fatalf("parseVersionParts dev suffix = %v, want [1 22]", got)
	}
}

func TestExpectedInstallEvents_matchesInstallCounts(t *testing.T) {
	for _, agent := range []string{"claude", "copilot", "cursor"} {
		events, err := ExpectedInstallEvents(agent)
		if err != nil {
			t.Fatalf("agent %s: %v", agent, err)
		}
		if len(events) == 0 {
			t.Fatalf("agent %s: no events", agent)
		}
	}
	claude, _ := ExpectedInstallEvents("claude")
	if len(claude) != len(agclaude.EventForKind) {
		t.Fatalf("claude events = %d, want %d", len(claude), len(agclaude.EventForKind))
	}
	copilot, _ := ExpectedInstallEvents("copilot")
	if len(copilot) != len(agcopilot.EventForKind) {
		t.Fatalf("copilot events = %d, want %d", len(copilot), len(agcopilot.EventForKind))
	}
	claudeAlias, err := ExpectedInstallEvents("claude-code")
	if err != nil {
		t.Fatalf("claude-code alias: %v", err)
	}
	if len(claudeAlias) != len(claude) {
		t.Fatalf("claude-code alias events = %d, want %d", len(claudeAlias), len(claude))
	}
	if _, err := ExpectedInstallEvents("nosuch"); err == nil {
		t.Fatal("expected error for unknown agent")
	}
}
