package dialect

import (
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "claude", in: "claude", want: sdkclaude.Dialect},
		{name: "claude-code alias", in: "claude-code", want: sdkclaude.Dialect},
		{name: "claudecode alias", in: "claudecode", want: sdkclaude.Dialect},
		{name: "copilot", in: "copilot", want: sdkcopilot.Dialect},
		{name: "github-copilot alias", in: "github-copilot", want: sdkcopilot.Dialect},
		{name: "gh alias", in: "gh", want: sdkcopilot.Dialect},
		{name: "cursor", in: "cursor", want: sdkcursor.Dialect},
		{name: "trimmed whitespace", in: "  CURSOR  ", want: sdkcursor.Dialect},
		{name: "unknown", in: "vscode", want: ""},
		{name: "empty", in: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.in); got != tt.want {
				t.Fatalf("Parse(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	if err := Validate(""); err != nil {
		t.Fatalf("empty: %v", err)
	}
	if err := Validate("claude"); err != nil {
		t.Fatalf("claude: %v", err)
	}
	if err := Validate("nosuch"); err == nil {
		t.Fatal("expected error for unknown agent")
	}
}
