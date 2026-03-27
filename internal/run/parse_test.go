package run

import (
	"slices"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestParseRunArgs_successDefaultFilePattern(t *testing.T) {
	mock := cli.NewMockConsole()
	args, pattern, err := parseRunArgs(mock, []string{"echo", "hi"})
	if err != nil {
		t.Fatalf("parseRunArgs: %v", err)
	}
	if pattern != defaultFilePatternFlagValue {
		t.Fatalf("file pattern: got %q want default %q", pattern, defaultFilePatternFlagValue)
	}
	want := []string{"echo", "hi"}
	if !slices.Equal(args, want) {
		t.Fatalf("args template: got %v want %v", args, want)
	}
}

func TestParseRunArgs_successWithFilePatternShorthand(t *testing.T) {
	mock := cli.NewMockConsole()
	args, pattern, err := parseRunArgs(mock, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("parseRunArgs: %v", err)
	}
	if pattern != `[.]go$` {
		t.Fatalf("file pattern: got %q", pattern)
	}
	want := []string{"echo", "x"}
	if !slices.Equal(args, want) {
		t.Fatalf("args template: got %v want %v", args, want)
	}
}

func TestParseRunArgs_successWithFilePatternLongForm(t *testing.T) {
	mock := cli.NewMockConsole()
	args, pattern, err := parseRunArgs(mock, []string{"--file-pattern", `[.]go$`, "echo", "y"})
	if err != nil {
		t.Fatalf("parseRunArgs: %v", err)
	}
	if pattern != `[.]go$` {
		t.Fatalf("file pattern: got %q", pattern)
	}
	if !slices.Equal(args, []string{"echo", "y"}) {
		t.Fatalf("args template: got %v", args)
	}
}

func TestParseRunArgs_lastFilePatternFlagWins(t *testing.T) {
	mock := cli.NewMockConsole()
	args, pattern, err := parseRunArgs(mock, []string{"-f", `first`, "-f", `[.]go$`, "echo", "z"})
	if err != nil {
		t.Fatalf("parseRunArgs: %v", err)
	}
	if pattern != `[.]go$` {
		t.Fatalf("want last -f value, got %q", pattern)
	}
	if !slices.Equal(args, []string{"echo", "z"}) {
		t.Fatalf("args: got %v", args)
	}
}

func TestParseRunArgs_nilProgramArgsMissingCommand(t *testing.T) {
	mock := cli.NewMockConsole()
	_, _, err := parseRunArgs(mock, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !mock.StderrContains("missing command to run") {
		t.Fatalf("stderr: %q", mock.StderrString())
	}
	if !mock.StderrContains("Usage:") {
		t.Fatalf("expected run help on stderr, got %q", mock.StderrString())
	}
}

func TestParseRunArgs_onlyFlagsMissingCommand(t *testing.T) {
	mock := cli.NewMockConsole()
	_, _, err := parseRunArgs(mock, []string{"-f", "*.go"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !mock.StderrContains("missing command to run") {
		t.Fatalf("stderr: %q", mock.StderrString())
	}
	if !mock.StderrContains("Usage:") {
		t.Fatalf("expected run help, got %q", mock.StderrString())
	}
}

func TestParseRunArgs_emptyFilePatternValue(t *testing.T) {
	mock := cli.NewMockConsole()
	_, _, err := parseRunArgs(mock, []string{"-f=", "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !mock.StderrContains("file-pattern value cannot be empty") {
		t.Fatalf("stderr: %q", mock.StderrString())
	}
	if !mock.StderrContains("Usage:") {
		t.Fatalf("expected run help, got %q", mock.StderrString())
	}
}

func TestParseRunArgs_unknownFlag(t *testing.T) {
	mock := cli.NewMockConsole()
	_, _, err := parseRunArgs(mock, []string{"-wat-not-a-flag", "echo"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !mock.StderrContains("Usage:") {
		t.Fatalf("expected run help after flag error, got %q", mock.StderrString())
	}
}

func TestCompileRunFilePattern_defaultMeansNoFilter(t *testing.T) {
	re, err := compileRunFilePattern(defaultFilePatternFlagValue)
	if err != nil {
		t.Fatalf("compileRunFilePattern: %v", err)
	}
	if re != nil {
		t.Fatal("expected nil regexp")
	}
}

func TestCompileRunFilePattern_valid(t *testing.T) {
	re, err := compileRunFilePattern(`\.go$`)
	if err != nil {
		t.Fatalf("compileRunFilePattern: %v", err)
	}
	if re == nil {
		t.Fatal("expected non-nil regexp")
	}
	if !re.MatchString("foo.go") || re.MatchString("foo.txt") {
		t.Fatalf("regexp match behavior wrong: MatchString foo.go=%v foo.txt=%v", re.MatchString("foo.go"), re.MatchString("foo.txt"))
	}
}

func TestCompileRunFilePattern_invalid(t *testing.T) {
	_, err := compileRunFilePattern(`(`)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid --file-pattern regexp") {
		t.Fatalf("error: %v", err)
	}
}
