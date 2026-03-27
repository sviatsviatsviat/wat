package run

import (
	"runtime"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestRunSubprocess_Success(t *testing.T) {
	args := subprocessSuccessArgs()
	mockConsole := cli.NewMockConsole()
	code := runSubprocess(mockConsole, args)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestRunSubprocess_FailureExitCode(t *testing.T) {
	args, expected := subprocessFailureArgs()
	mockConsole := cli.NewMockConsole()
	code := runSubprocess(mockConsole, args)
	if code != expected {
		t.Fatalf("expected exit %d, got %d", expected, code)
	}
}

func TestRunSubprocess_NoCommandAfterTemplating(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	code := runSubprocess(mockConsole, []string{})
	if code != cli.ExitBadInput {
		t.Fatalf("expected exit cli.ExitBadInput, got %d, stderr=%q", code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("no command to execute after templating") {
		t.Fatalf("expected templating message on stderr, got %q", mockConsole.StderrString())
	}
}

func subprocessSuccessArgs() []string {
	return []string{"echo", "ok"}
}

func subprocessFailureArgs() ([]string, int) {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/C", "exit 7"}, 7
	}
	return []string{"sh", "-c", "exit 7"}, 7
}
