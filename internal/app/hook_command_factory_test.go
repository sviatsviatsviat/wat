package app

import (
	"errors"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

func TestNewHookCommand_Run(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := newHookCommand("run", mockConsole, runner, []string{"echo", "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hookCommand == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestNewHookCommand_RunWithFilePatternFlag(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := newHookCommand("run", mockConsole, runner, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hookCommand == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestNewHookCommand_RunEmptyArgv(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := newHookCommand("run", mockConsole, runner, []string{})
	if err == nil {
		t.Fatal("expected error from run.NewRunCommand")
	}
	if hookCommand != nil {
		t.Fatal("expected nil command")
	}
	if errors.Is(err, errHookCommandBadInput) {
		t.Fatal("run parse error should not be errHookCommandBadInput")
	}
	if !mockConsole.StderrContains("missing command to run") {
		t.Fatalf("expected run help on stderr, got %q", mockConsole.StderrString())
	}
}

func TestNewHookCommand_UnknownSubcommand(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := newHookCommand("nope", mockConsole, nil, []string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if hookCommand != nil {
		t.Fatalf("expected nil command, got %#v", hookCommand)
	}
	if !errors.Is(err, errHookCommandBadInput) {
		t.Fatalf("expected errHookCommandBadInput, got %v", err)
	}
	if !mockConsole.StderrContains(`unknown command "nope"`) {
		t.Fatalf("expected unknown command on stderr, got %q", mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("wat <host>") {
		t.Fatalf("expected root help on stderr, got %q", mockConsole.StderrString())
	}
}
