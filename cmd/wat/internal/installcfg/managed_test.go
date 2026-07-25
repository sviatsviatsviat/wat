package installcfg

import "testing"

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

func TestIsWatManagedAgentCommand_ignoresEvent(t *testing.T) {
	watAbs := "/usr/bin/wat"
	if !IsWatManagedAgentCommand(watAbs+" run --agent cursor --event beforeShellExecution", "cursor", watAbs) {
		t.Fatal("expected cursor command to match regardless of event")
	}
	if IsWatManagedAgentCommand(watAbs+" run --agent claude --event PreToolUse", "cursor", watAbs) {
		t.Fatal("command for another agent must not match")
	}
}
