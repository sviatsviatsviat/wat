package app

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func validCursorHookJSON() string {
	return `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "c1",
		"generation_id": "g1",
		"model": "claude-test",
		"cursor_version": "1.0",
		"workspace_roots": ["/a"],
		"user_email": "a@b.com",
		"transcript_path": "/t"
	}`
}

func runEchoHookEventArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"run", "cmd", "/C", "echo __HOOK_EVENT_NAME__ 1>&2"}
	}
	return []string{"run", "sh", "-c", "echo __HOOK_EVENT_NAME__ >&2"}
}

func goVersionToStderrArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"run", "cmd", "/C", "go version 1>&2"}
	}
	return []string{"run", "sh", "-c", "go version >&2"}
}

func assertHookStdoutJSON(t *testing.T, stdout string) {
	t.Helper()
	if strings.TrimSpace(stdout) != "{}" {
		t.Fatalf("expected wat stdout JSON \"{}\", got %q", stdout)
	}
}

func assertStdoutEmpty(t *testing.T, stdout string) {
	t.Helper()
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

func TestExecute_RunRendersTemplateAndExecutes(t *testing.T) {
	stdin := strings.NewReader(validCursorHookJSON())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Execute(runEchoHookEventArgs(), stdin, &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
	if !strings.Contains(stderr.String(), "afterFileEdit") {
		t.Fatalf("expected child echo on stderr, got %q", stderr.String())
	}
}

func TestExecute_RunRunsPlainCommand(t *testing.T) {
	stdin := strings.NewReader(validCursorHookJSON())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Execute(goVersionToStderrArgs(), stdin, &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
	if !strings.Contains(stderr.String(), "go version") {
		t.Fatalf("expected go version on stderr, got %q", stderr.String())
	}
}

func TestExecute_RunExplicitHostCursorStillWorks(t *testing.T) {
	stdin := strings.NewReader(validCursorHookJSON())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := append([]string{"run", "-H", "cursor"}, goVersionToStderrArgs()[1:]...)
	code := Execute(args, stdin, &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
	if !strings.Contains(stderr.String(), "go version") {
		t.Fatalf("expected go version on stderr, got %q", stderr.String())
	}
}

func TestExecute_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "wat <command>") {
		t.Fatalf("expected root help output, got %q", stderr.String())
	}
}

func TestExecute_InvalidJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(goVersionToStderrArgs(), strings.NewReader("{bad"), &stdout, &stderr)
	if code != cli.ExitGeneral {
		t.Fatalf("expected cli.ExitGeneral for stdin JSON parse error, got %d", code)
	}
	if !strings.Contains(stderr.String(), "failed to parse stdin event JSON") {
		t.Fatalf("expected parse error message, got %q", stderr.String())
	}
	assertStdoutEmpty(t, stdout.String())
}

func TestExecute_RunHelpWhenNoCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"run"}, strings.NewReader("{}"), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing command to run") {
		t.Fatalf("expected missing command message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "__CONVERSATION_ID__") {
		t.Fatalf("expected run help with placeholder list, got %q", stderr.String())
	}
}

func TestExecute_DefaultHostIsCursor(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(goVersionToStderrArgs(), strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess without --host, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
}

func TestExecute_MissingCommandAfterHost(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"run", "--host", "cursor"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing command to run") {
		t.Fatalf("expected missing command message, got %q", stderr.String())
	}
}

func TestExecute_UnsupportedHost(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	args := append([]string{"run", "--host", "other"}, goVersionToStderrArgs()[1:]...)
	code := Execute(args, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput for unsupported host, got %d", code)
	}
	if !strings.Contains(stderr.String(), "not supported yet") {
		t.Fatalf("expected unsupported host message, got %q", stderr.String())
	}
	assertStdoutEmpty(t, stdout.String())
}

func TestExecute_UnsupportedCursorEvent(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(goVersionToStderrArgs(), strings.NewReader(`{"hook_event_name":"preToolUse"}`), &stdout, &stderr)
	if code != cli.ExitGeneral {
		t.Fatalf("expected cli.ExitGeneral for unsupported cursor event, got %d", code)
	}
	if !strings.Contains(stderr.String(), "not supported yet") {
		t.Fatalf("expected unsupported event message, got %q", stderr.String())
	}
	assertStdoutEmpty(t, stdout.String())
}

func TestExecute_UnknownTemplatePlaceholder(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var args []string
	if runtime.GOOS == "windows" {
		args = []string{"run", "cmd", "/C", "echo __FILE__ 1>&2"}
	} else {
		args = []string{"run", "sh", "-c", "echo __FILE__ >&2"}
	}
	code := Execute(args, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d, stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "unknown template placeholders") || !strings.Contains(stderr.String(), "FILE") {
		t.Fatalf("expected unknown placeholder error, got %q", stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
}

func TestExecute_UnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"nope"}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("expected unknown command message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "wat <command>") {
		t.Fatalf("expected root help after unknown command, got %q", stderr.String())
	}
}

func TestExecute_ParseErrorMissingHostValue(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"run", "--host"}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing host value after --host") {
		t.Fatalf("expected missing host value message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "wat <command>") {
		t.Fatalf("expected root help after parse error, got %q", stderr.String())
	}
}

func TestExecute_ParseErrorEmptyHostValue(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"run", "--host", ""}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "host value cannot be empty") {
		t.Fatalf("expected empty host message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "wat <command>") {
		t.Fatalf("expected root help after parse error, got %q", stderr.String())
	}
}

func TestExecute_ParseErrorMissingHostValueDashCapitalH(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"run", "-H"}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing host value after -H") {
		t.Fatalf("expected missing host value message, got %q", stderr.String())
	}
}
