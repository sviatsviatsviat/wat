package run

import (
	"runtime"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/watexec"
)

func TestNewRunCommand_EmptyArgv(t *testing.T) {
	tests := []struct {
		name string
		argv []string
	}{
		{name: "nil", argv: nil},
		{name: "empty", argv: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConsole := cli.NewMockConsole()
			runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
			hookCommand, err := NewRunCommand(mockConsole, runner, tt.argv)
			if err == nil {
				t.Fatal("expected error")
			}
			if hookCommand != nil {
				t.Fatal("expected nil command")
			}
			if !mockConsole.StderrContains("missing command to run") {
				t.Fatalf("stderr missing error line, got %q", mockConsole.StderrString())
			}
			if !mockConsole.StderrContains("Usage:") {
				t.Fatalf("stderr missing run help, got %q", mockConsole.StderrString())
			}
		})
	}
}

func TestNewRunCommand_OK(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hookCommand == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestRunCommand_Execute_NilHookContext(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	code := hookCommand.Execute(nil)
	if code != cli.ExitGeneral {
		t.Fatalf("expected ExitGeneral, got %d, stderr=%q", code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("HookContext is nil") {
		t.Fatalf("expected nil context error on stderr, got %q", mockConsole.StderrString())
	}
}

func TestRunCommand_Execute_UnknownPlaceholder(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"echo", "__NOT_A_SUPPORTED_KEY__"})
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	ctx := testHookContext(stubTemplateBindings{defined: map[string]struct{}{}})
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

func TestRunCommand_Execute_SubstitutionAndSuccess(t *testing.T) {
	var argv []string
	if runtime.GOOS == "windows" {
		argv = []string{"cmd", "/C", "echo __CONVERSATION_ID__"}
	} else {
		argv = []string{"sh", "-c", "echo __CONVERSATION_ID__"}
	}
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, argv)
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	bindings := stubTemplateBindings{
		defined: map[string]struct{}{"CONVERSATION_ID": {}},
		values:  map[string]string{"CONVERSATION_ID": "conv-test-1"},
	}
	ctx := testHookContext(bindings)
	code := hookCommand.Execute(ctx)
	if code != cli.ExitSuccess {
		t.Fatalf("expected success, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestRunCommand_Execute_SubprocessFailureExitCode(t *testing.T) {
	var argv []string
	if runtime.GOOS == "windows" {
		argv = []string{"cmd", "/C", "exit 9"}
	} else {
		argv = []string{"sh", "-c", "exit 9"}
	}
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, argv)
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	ctx := testHookContext(stubTemplateBindings{defined: map[string]struct{}{}})
	code := hookCommand.Execute(ctx)
	if code != 9 {
		t.Fatalf("expected subprocess exit 9, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

type stubTemplateBindings struct {
	defined map[string]struct{}
	values  map[string]string
}

func (stub stubTemplateBindings) TemplateValue(key string) (string, bool) {
	if stub.defined == nil {
		return "", false
	}
	if _, ok := stub.defined[key]; !ok {
		return "", false
	}
	if stub.values == nil {
		return "", true
	}
	return stub.values[key], true
}

func testHookContext(bindings core.TemplateBindings) *core.HookContext {
	return &core.HookContext{
		TemplateBindings: bindings,
	}
}

func TestNewRunCommand_InvalidFilePatternRegexp(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	_, err := NewRunCommand(mockConsole, runner, []string{"-f", `(`, "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid --file-pattern regexp") {
		t.Fatalf("error: %v", err)
	}
}

func TestRunCommand_Execute_FilePatternNoMatchSkipsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	bindings := stubTemplateBindings{
		defined: map[string]struct{}{"FILE_PATH": {}},
		values:  map[string]string{"FILE_PATH": `D:\repo\file.txt`},
	}
	code := hookCommand.Execute(testHookContext(bindings))
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess when path does not match, got %d", code)
	}
}

func TestRunCommand_Execute_FilePatternMatchRunsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	bindings := stubTemplateBindings{
		defined: map[string]struct{}{"FILE_PATH": {}},
		values:  map[string]string{"FILE_PATH": `D:\repo\file.go`},
	}
	code := hookCommand.Execute(testHookContext(bindings))
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess from echo, got %d, stderr=%q", code, mockConsole.StderrString())
	}
}

func TestRunCommand_Execute_FilePatternIgnoredWithoutFilePathBinding(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	runner := watexec.NewRunner(mockConsole.StderrBufferWriter(), mockConsole)
	hookCommand, err := NewRunCommand(mockConsole, runner, []string{"-f", `[.]go$`, "echo", "y"})
	if err != nil {
		t.Fatalf("NewRunCommand: %v", err)
	}
	code := hookCommand.Execute(testHookContext(stubTemplateBindings{defined: map[string]struct{}{}}))
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess, got %d", code)
	}
}
