package execcommand

import (
	"runtime"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func TestNewExecCommand_EmptyArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "nil", args: nil},
		{name: "empty", args: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConsole := cli.NewMockConsole()
			hookCommand, err := NewExecCommand(mockConsole, tt.args)
			if err == nil {
				t.Fatal("expected error")
			}
			if hookCommand != nil {
				t.Fatal("expected nil command")
			}
			if !mockConsole.StderrContains("missing subprocess command") {
				t.Fatalf("stderr missing error line, got %q", mockConsole.StderrString())
			}
			if !mockConsole.StderrContains("Usage:") {
				t.Fatalf("stderr missing exec help, got %q", mockConsole.StderrString())
			}
		})
	}
}

func TestNewExecCommand_OK(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hookCommand == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestExecCommand_Execute_NilHookContext(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	code := hookCommand.Execute(nil)
	if code != cli.ExitGeneral {
		t.Fatalf("expected ExitGeneral, got %d, stderr=%q", code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("HookContext is nil") {
		t.Fatalf("expected nil context error on stderr, got %q", mockConsole.StderrString())
	}
}

func TestExecCommand_Execute_UnknownPlaceholder(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"echo", "__NOT_A_SUPPORTED_KEY__"})
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	ctx := testExecHookContext(cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	code := hookCommand.Execute(ctx)
	if code != cli.ExitBadInput {
		t.Fatalf("expected ExitBadInput, got %d, stderr=%q", code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("unknown template placeholders") {
		t.Fatalf("expected placeholder error on stderr, got %q", mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("NOT_A_SUPPORTED_KEY") {
		t.Fatalf("expected unknown key name on stderr, got %q", mockConsole.StderrString())
	}
}

func TestExecCommand_Execute_SubstitutionAndSuccess(t *testing.T) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", "echo __CONVERSATION_ID__"}
	} else {
		cmdArgs = []string{"sh", "-c", "echo __CONVERSATION_ID__"}
	}
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, cmdArgs)
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	ctx := testExecHookContext(cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{
			HookEventName:  "sessionEnd",
			ConversationID: "conv-test-1",
		},
	})
	code := hookCommand.Execute(ctx)
	if code != cli.ExitSuccess {
		t.Fatalf("expected success, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestExecCommand_Execute_SubprocessFailureExitCode(t *testing.T) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", "exit 9"}
	} else {
		cmdArgs = []string{"sh", "-c", "exit 9"}
	}
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, cmdArgs)
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	ctx := testExecHookContext(cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	code := hookCommand.Execute(ctx)
	if code != 9 {
		t.Fatalf("expected subprocess exit 9, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func testExecHookContext[T any](data cursor.CursorHookRunData[T]) *core.HookContext {
	return &core.HookContext{
		HookHost:   cursor.HookHostCursor,
		ParsedData: &data,
	}
}

func TestNewExecCommand_InvalidFilePatternRegexp(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	_, err := NewExecCommand(mockConsole, []string{"-f", `(`, "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid --file-pattern regexp") {
		t.Fatalf("error: %v", err)
	}
}

func TestExecCommand_Execute_FilePatternNoMatchSkipsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{
			FilePath: `D:\repo\file.txt`,
		},
	}
	code := hookCommand.Execute(testExecHookContext(data))
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess when path does not match, got %d", code)
	}
}

func TestExecCommand_Execute_FilePatternMatchRunsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{
			FilePath: `D:\repo\file.go`,
		},
	}
	code := hookCommand.Execute(testExecHookContext(data))
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess from echo, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestExecCommand_Execute_FilePatternIgnoredWithoutFilePathBinding(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	hookCommand, err := NewExecCommand(mockConsole, []string{"-f", `[.]go$`, "echo", "y"})
	if err != nil {
		t.Fatalf("NewExecCommand: %v", err)
	}
	ctx := testExecHookContext(cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	code := hookCommand.Execute(ctx)
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess, got %d", code)
	}
}
