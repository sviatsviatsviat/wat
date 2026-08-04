package installcfg

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseWatRunFlags(t *testing.T) {
	agent, event, ok := ParseWatRunFlags("/usr/bin/wat run --agent claude --event PreToolUse")
	if !ok || agent != "claude" || event != "PreToolUse" {
		t.Fatalf("got %q %q %v", agent, event, ok)
	}
}

func TestWatRunCommand_quotesPathsWithSpaces(t *testing.T) {
	watAbs := filepath.Join("C:", "Program Files", "wat", "wat.exe")
	got := WatRunCommand(watAbs, "claude", "PreToolUse")
	if !strings.HasPrefix(got, `"`) || !strings.Contains(got, `" run --agent claude --event PreToolUse`) {
		t.Fatalf("WatRunCommand = %q", got)
	}
	if !IsWatManagedCommand(got, "claude", "PreToolUse", watAbs) {
		t.Fatalf("quoted command not recognized as managed: %q", got)
	}
	plain := WatRunCommand(`/usr/bin/wat`, "cursor", "sessionStart")
	if plain != `/usr/bin/wat run --agent cursor --event sessionStart` {
		t.Fatalf("plain WatRunCommand = %q", plain)
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

func TestIsWatManagedCommand_recognizesQuotedAndBasename(t *testing.T) {
	watAbs := `/usr/local/bin/wat`
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "quoted absolute path",
			command: `"/usr/local/bin/wat" run --agent claude --event PreToolUse`,
			want:    true,
		},
		{
			name:    "quoted path with spaces",
			command: `"C:\Program Files\wat\wat.exe" run --agent claude --event PreToolUse`,
			want:    true,
		},
		{
			name:    "basename with exe suffix",
			command: `Wat.EXE run --agent claude --event PreToolUse`,
			want:    true,
		},
		{
			name:    "other path with wat basename",
			command: `/other/wat run --agent claude --event PreToolUse`,
			want:    true,
		},
		{
			name:    "non-wat basename",
			command: `/usr/bin/hook run --agent claude --event PreToolUse`,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWatManagedCommand(tt.command, "claude", "PreToolUse", watAbs)
			if got != tt.want {
				t.Fatalf("IsWatManagedCommand(%q) = %v, want %v", tt.command, got, tt.want)
			}
			gotAgent := IsWatManagedAgentCommand(tt.command, "claude", watAbs)
			if gotAgent != tt.want {
				t.Fatalf("IsWatManagedAgentCommand(%q) = %v, want %v", tt.command, gotAgent, tt.want)
			}
		})
	}
}
