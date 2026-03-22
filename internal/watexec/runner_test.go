package watexec

import (
	"runtime"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestRunner_Run_Success(t *testing.T) {
	args := successCommand()
	mockConsole := cli.NewMockConsole()
	runner := NewRunner(mockConsole.StderrBufferWriter(), mockConsole)

	code := runner.Run(args)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestRunner_Run_FailureExitCode(t *testing.T) {
	args, expected := failureCommand()
	mockConsole := cli.NewMockConsole()
	runner := NewRunner(mockConsole.StderrBufferWriter(), mockConsole)

	code := runner.Run(args)
	if code != expected {
		t.Fatalf("expected exit %d, got %d", expected, code)
	}
}

func TestRunner_Run_NoCommandAfterTemplating(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	code := runner.Run([]string{})
	if code != cli.ExitBadInput {
		t.Fatalf("expected exit cli.ExitBadInput, got %d, stderr=%q", code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("no command to execute after templating") {
		t.Fatalf("expected templating message on stderr, got %q", mockConsole.StderrString())
	}
}

func successCommand() []string {
	// echo is a Windows shell builtin; on Unix it is usually a real binary on PATH.
	return []string{"echo", "ok"}
}

func failureCommand() ([]string, int) {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/C", "exit 7"}, 7
	}
	return []string{"sh", "-c", "exit 7"}, 7
}
