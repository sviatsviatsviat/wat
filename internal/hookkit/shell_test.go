package hookkit

import (
	"encoding/json"
	"testing"
)

func TestExtractShellCommand(t *testing.T) {
	t.Parallel()
	if got := ExtractShellCommand(json.RawMessage(`{"command":"echo hi"}`)); got != "echo hi" {
		t.Fatalf("ExtractShellCommand() = %q", got)
	}
	if got := ExtractShellCommand(nil); got != "" {
		t.Fatalf("ExtractShellCommand(nil) = %q", got)
	}
}

func TestShellSingleQuote(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"":        "''",
		"hello":   "'hello'",
		"it's ok": "'it'\\''s ok'",
	}
	for in, want := range cases {
		if got := ShellSingleQuote(in); got != want {
			t.Fatalf("ShellSingleQuote(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidEnvKey(t *testing.T) {
	t.Parallel()
	if !ValidEnvKey("FOO_1") {
		t.Fatal("FOO_1 should be valid")
	}
	if ValidEnvKey("1BAD") {
		t.Fatal("1BAD should be invalid")
	}
}
