package installcfg

import (
	"slices"
	"testing"
)

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

func TestExpectedEvents(t *testing.T) {
	// Sorted wire names ExpectedEvents must return (independent of installcfg event slices).
	wantClaude := []string{
		"Notification",
		"PermissionRequest",
		"PostToolUse",
		"PostToolUseFailure",
		"PreCompact",
		"PreToolUse",
		"SessionEnd",
		"SessionStart",
		"Stop",
		"StopFailure",
		"SubagentStart",
		"SubagentStop",
		"UserPromptSubmit",
	}
	wantCopilot := []string{
		"ErrorOccurred",
		"Notification",
		"PermissionRequest",
		"PostToolUse",
		"PostToolUseFailure",
		"PreCompact",
		"PreToolUse",
		"SessionEnd",
		"SessionStart",
		"Stop",
		"SubagentStart",
		"SubagentStop",
		"UserPromptSubmit",
	}
	wantCursor := []string{
		"afterFileEdit",
		"afterMCPExecution",
		"afterShellExecution",
		"beforeMCPExecution",
		"beforeReadFile",
		"beforeShellExecution",
		"beforeSubmitPrompt",
		"postToolUse",
		"postToolUseFailure",
		"preCompact",
		"preToolUse",
		"sessionEnd",
		"sessionStart",
		"stop",
		"subagentStart",
		"subagentStop",
	}

	tests := []struct {
		agent string
		want  []string
	}{
		{agent: "claude", want: wantClaude},
		{agent: "copilot", want: wantCopilot},
		{agent: "cursor", want: wantCursor},
		{agent: "claude-code", want: wantClaude},
	}
	for _, tt := range tests {
		t.Run(tt.agent, func(t *testing.T) {
			got, err := ExpectedEvents(tt.agent)
			if err != nil {
				t.Fatalf("ExpectedEvents(%q): %v", tt.agent, err)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("ExpectedEvents(%q) = %v, want %v", tt.agent, got, tt.want)
			}
		})
	}

	if _, err := ExpectedEvents("nosuch"); err == nil {
		t.Fatal("expected error for unknown agent")
	}
}

func TestIsValidEvent(t *testing.T) {
	if !IsValidEvent("claude", "PreToolUse") {
		t.Fatal("expected PreToolUse valid for claude")
	}
	if !IsValidEvent("cursor", "beforeShellExecution") {
		t.Fatal("expected Cursor dedicated event valid")
	}
	if IsValidEvent("claude", "beforeShellExecution") {
		t.Fatal("Claude must not accept Cursor dedicated events")
	}
	if IsValidEvent("nosuch", "PreToolUse") {
		t.Fatal("unknown agent must be invalid")
	}
}
