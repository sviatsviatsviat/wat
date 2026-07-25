package paths

import (
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func TestConfigPath(t *testing.T) {
	root := "/project"
	tests := []struct {
		dialect string
		want    string
	}{
		{dialect: sdkclaude.Dialect, want: "/project/.claude/settings.json"},
		{dialect: sdkcopilot.Dialect, want: "/project/.github/hooks/wat.json"},
		{dialect: sdkcursor.Dialect, want: "/project/.cursor/hooks.json"},
		{dialect: "unknown", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.dialect, func(t *testing.T) {
			if got := ConfigPath(tt.dialect, root); got != tt.want {
				t.Fatalf("ConfigPath(%q) = %q, want %q", tt.dialect, got, tt.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	root := "/project"
	got := All(root)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	want := []AgentConfig{
		{Agent: sdkclaude.Dialect, Path: ConfigPath(sdkclaude.Dialect, root)},
		{Agent: sdkcopilot.Dialect, Path: ConfigPath(sdkcopilot.Dialect, root)},
		{Agent: sdkcursor.Dialect, Path: ConfigPath(sdkcursor.Dialect, root)},
	}
	for i, cfg := range want {
		if got[i].Agent != cfg.Agent {
			t.Fatalf("%s agent = %q, want %q", cfg.Agent, got[i].Agent, cfg.Agent)
		}
		if got[i].Path != cfg.Path {
			t.Fatalf("%s path = %q, want %q", cfg.Agent, got[i].Path, cfg.Path)
		}
	}
}
