package execcommand

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func TestNewExecHookHandlerProvider_EmptyArgs(t *testing.T) {
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
			provider, err := NewExecHookHandlerProvider(mockConsole, tt.args)
			if err == nil {
				t.Fatal("expected error")
			}
			if provider != nil {
				t.Fatal("expected nil provider")
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

func TestNewExecHookHandlerProvider_OK(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestExecHookHandler_Handle_UnknownPlaceholder(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"echo", "__NOT_A_SUPPORTED_KEY__"})
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	adapter := testExecHookAdapter(mockConsole, cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	handler, hookErr := provider.HookHandlerFor(adapter)
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitBadInput {
		t.Fatalf("expected ExitBadInput, got %d, stderr=%q", result.Code, mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("unknown template placeholders") {
		t.Fatalf("expected placeholder error on stderr, got %q", mockConsole.StderrString())
	}
	if !mockConsole.StderrContains("NOT_A_SUPPORTED_KEY") {
		t.Fatalf("expected unknown key name on stderr, got %q", mockConsole.StderrString())
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func TestExecHookHandler_Handle_AfterAgentThought_substitutesDurationMs(t *testing.T) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", "echo __DURATION_MS__ 1>&2"}
	} else {
		cmdArgs = []string{"sh", "-c", "echo __DURATION_MS__ >&2"}
	}
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, cmdArgs)
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	adapter := testExecHookAdapter(mockConsole, cursor.CursorHookRunData[cursor.AfterAgentThoughtFields]{
		Common: cursor.HookDataCommon{
			HookEventName:  "afterAgentThought",
			ConversationID: "conv-ms",
		},
		EventSpecific: &cursor.AfterAgentThoughtFields{
			Text:       "reasoning",
			DurationMs: 5000,
		},
	})
	handler, hookErr := provider.HookHandlerFor(adapter)
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitSuccess {
		t.Fatalf("expected success, got %d, stderr=%q", result.Code, mockConsole.StderrString())
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
	if !strings.Contains(mockConsole.StderrString(), "5000") {
		t.Fatalf("expected duration in child stderr, got %q", mockConsole.StderrString())
	}
}

func TestExecHookHandler_Handle_SubstitutionAndSuccess(t *testing.T) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", "echo __CONVERSATION_ID__"}
	} else {
		cmdArgs = []string{"sh", "-c", "echo __CONVERSATION_ID__"}
	}
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, cmdArgs)
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	adapter := testExecHookAdapter(mockConsole, cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{
			HookEventName:  "sessionEnd",
			ConversationID: "conv-test-1",
		},
	})
	handler, hookErr := provider.HookHandlerFor(adapter)
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitSuccess {
		t.Fatalf("expected success, got %d, stderr=%q", result.Code, mockConsole.StderrString())
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func TestExecHookHandler_Handle_SubprocessFailureExitCode(t *testing.T) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd", "/C", "exit 9"}
	} else {
		cmdArgs = []string{"sh", "-c", "exit 9"}
	}
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, cmdArgs)
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	adapter := testExecHookAdapter(mockConsole, cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	handler, hookErr := provider.HookHandlerFor(adapter)
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != 9 {
		t.Fatalf("expected subprocess exit 9, got %d, stderr=%q", result.Code, mockConsole.StderrString())
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func testExecHookAdapter[T any](mock *cli.MockConsole, data cursor.CursorHookRunData[T]) core.HookAdapter {
	return cursor.NewHookAdapter(mock, data.Common, data.EventSpecific)
}

func TestNewExecHookHandlerProvider_InvalidFilePatternRegexp(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	_, err := NewExecHookHandlerProvider(mockConsole, []string{"-f", `(`, "echo", "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid --file-pattern regexp") {
		t.Fatalf("error: %v", err)
	}
}

func TestExecHookHandler_Handle_FilePatternNoMatchSkipsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{
			FilePath: `D:\repo\file.txt`,
		},
	}
	handler, hookErr := provider.HookHandlerFor(testExecHookAdapter(mockConsole, data))
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess when path does not match, got %d", result.Code)
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func TestExecHookHandler_Handle_FilePatternMatchRunsSubprocess(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	data := cursor.CursorHookRunData[cursor.AfterFileEditFields]{
		Common: cursor.HookDataCommon{HookEventName: "afterFileEdit"},
		EventSpecific: &cursor.AfterFileEditFields{
			FilePath: `D:\repo\file.go`,
		},
	}
	handler, hookErr := provider.HookHandlerFor(testExecHookAdapter(mockConsole, data))
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess from echo, got %d, stderr=%q", result.Code, mockConsole.StderrString())
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func TestExecHookHandler_Handle_FilePatternIgnoredWithoutFilePathBinding(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"-f", `[.]go$`, "echo", "y"})
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	adapter := testExecHookAdapter(mockConsole, cursor.CursorHookRunData[struct{}]{
		Common: cursor.HookDataCommon{HookEventName: "sessionEnd"},
	})
	handler, hookErr := provider.HookHandlerFor(adapter)
	if hookErr != nil {
		t.Fatalf("HookHandlerFor: %v", hookErr)
	}
	result := handler.Handle()
	if result.Code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess, got %d", result.Code)
	}
	if mockConsole.StdoutString() != "{}\n" {
		t.Fatalf("hook stdout: want %q, got %q", "{}\n", mockConsole.StdoutString())
	}
}

func TestHookHandlerFor_unsupportedAdapterType(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := NewExecHookHandlerProvider(mockConsole, []string{"echo", "x"})
	if err != nil {
		t.Fatalf("NewExecHookHandlerProvider: %v", err)
	}
	_, hookErr := provider.HookHandlerFor(unsupportedHookAdapterStub{})
	if hookErr == nil {
		t.Fatal("expected error for unsupported adapter type")
	}
	if !errors.Is(hookErr, core.ErrHookAdapterNotSupported) {
		t.Fatalf("want %v, got %v", core.ErrHookAdapterNotSupported, hookErr)
	}
}

// unsupportedHookAdapterStub implements [core.HookAdapter] but is not a supported Cursor exec payload type.
type unsupportedHookAdapterStub struct{}

func (unsupportedHookAdapterStub) HookHost() string { return "cursor" }

func (unsupportedHookAdapterStub) ReturnEmpty() {}
