package run

import (
	"strings"
	"testing"
)

func TestCommandForArgs_emptyArgs(t *testing.T) {
	_, err := commandForArgs(nil)
	if err == nil {
		t.Fatal("expected error for nil args")
	}
	if !strings.Contains(err.Error(), "no command arguments provided") {
		t.Fatalf("error %q should mention missing arguments", err.Error())
	}
	_, err = commandForArgs([]string{})
	if err == nil {
		t.Fatal("expected error for empty args")
	}
	if !strings.Contains(err.Error(), "no command arguments provided") {
		t.Fatalf("error %q should mention missing arguments", err.Error())
	}
}

func TestWindowsShellArg_quoting(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: `""`},
		{name: "hello", in: "hello", want: "hello"},
		{name: "ampersand unquoted", in: "a&b", want: "a^&b"},
		{name: "pipe unquoted", in: "a|b", want: "a^|b"},
		{name: "gt unquoted", in: "a>b", want: "a^>b"},
		{name: "lt unquoted", in: "a<b", want: "a^<b"},
		{name: "percent unquoted", in: "50%", want: "50^%"},
		{name: "parens unquoted", in: "a(b)c", want: "a^(b^)c"},
		{name: "caret", in: "caret^", want: "caret^^"},
		{name: "space needs quotes", in: "a b", want: `"a b"`},
		{name: "embedded double quotes", in: `say "hi"`, want: `"say ""hi"""`},
		{name: "ampersand inside quoted", in: "a&b c", want: `"a&b c"`},
		{name: "percent inside quoted", in: "a 50%", want: `"a 50%%"`},
		{name: "bang inside quoted", in: "ok !x", want: `"ok ^!x"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := windowsShellArg(tt.in)
			if got != tt.want {
				t.Fatalf("windowsShellArg(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPosixShellLine_quoting(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "empty args",
			args: []string{},
			want: "",
		},
		{
			name: "empty single arg",
			args: []string{""},
			want: "''",
		},
		{
			name: "space in arg",
			args: []string{"a b"},
			want: "'a b'",
		},
		{
			name: "multiple single quotes",
			args: []string{`o'c'lock`},
			want: `'o'\''c'\''lock'`,
		},
		{
			name: "backslashes literal under single quotes",
			args: []string{`a\b`},
			want: `'a\b'`,
		},
		{
			name: "double quotes literal (not shell-escaped)",
			args: []string{`"hello"`},
			want: `'"hello"'`,
		},
		{
			name: "echo its ok",
			args: []string{"echo", "it's", "ok"},
			want: `'echo' 'it'\''s' 'ok'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := posixShellLine(tt.args)
			if got != tt.want {
				t.Fatalf("posixShellLine(%#v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
