package installcfg

import (
	"slices"
	"sort"
	"testing"

	portclaude "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/claude"
	portcopilot "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/copilot"
	portcursor "github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/cursor"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
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

func TestExpectedEvents_matchesInstallCounts(t *testing.T) {
	for _, agent := range []string{"claude", "copilot", "cursor"} {
		events, err := ExpectedEvents(agent)
		if err != nil {
			t.Fatalf("agent %s: %v", agent, err)
		}
		if len(events) == 0 {
			t.Fatalf("agent %s: no events", agent)
		}
	}

	claude, err := ExpectedEvents("claude")
	if err != nil {
		t.Fatalf("claude: %v", err)
	}
	wantClaude := sortedKindEvents(portclaude.EventForKind)
	if !slices.Equal(claude, wantClaude) {
		t.Fatalf("claude events = %v, want %v", claude, wantClaude)
	}

	copilot, err := ExpectedEvents("copilot")
	if err != nil {
		t.Fatalf("copilot: %v", err)
	}
	wantCopilot := sortedKindEvents(portcopilot.EventForKind)
	if !slices.Equal(copilot, wantCopilot) {
		t.Fatalf("copilot events = %v, want %v", copilot, wantCopilot)
	}

	cursor, err := ExpectedEvents("cursor")
	if err != nil {
		t.Fatalf("cursor: %v", err)
	}
	wantCursor := sortedKindEvents(portcursor.EventForKind)
	for ev := range portcursor.DedicatedEvents {
		wantCursor = append(wantCursor, ev)
	}
	sort.Strings(wantCursor)
	if !slices.Equal(cursor, wantCursor) {
		t.Fatalf("cursor events = %v, want %v", cursor, wantCursor)
	}

	claudeAlias, err := ExpectedEvents("claude-code")
	if err != nil {
		t.Fatalf("claude-code alias: %v", err)
	}
	if !slices.Equal(claudeAlias, claude) {
		t.Fatalf("claude-code alias events = %v, want %v", claudeAlias, claude)
	}
	if _, err := ExpectedEvents("nosuch"); err == nil {
		t.Fatal("expected error for unknown agent")
	}
}

func sortedKindEvents(m map[model.Kind]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		if v != "" {
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
