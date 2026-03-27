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
		"transcript_path": "/t",
		"file_path": "/repo/main.go",
		"edits": [{"old_string": "a", "new_string": "b"}]
	}`
}

func runEchoHookEventArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"cursor", "run", "cmd", "/C", "echo __HOOK_EVENT_NAME__ 1>&2"}
	}
	return []string{"cursor", "run", "sh", "-c", "echo __HOOK_EVENT_NAME__ >&2"}
}

func goVersionToStderrArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"cursor", "run", "cmd", "/C", "go version 1>&2"}
	}
	return []string{"cursor", "run", "sh", "-c", "go version >&2"}
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

func TestExecute_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "wat <host>") {
		t.Fatalf("expected root help output, got %q", stderr.String())
	}
}

func TestExecute_TooFewArgs(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"cursor"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "host") && !strings.Contains(stderr.String(), "command") {
		t.Fatalf("expected host/command message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "wat <host>") {
		t.Fatalf("expected root help, got %q", stderr.String())
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
	code := Execute([]string{"cursor", "run"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
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

func TestExecute_CursorHostRuns(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(goVersionToStderrArgs(), strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
}

func TestExecute_UnsupportedHost(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	args := append([]string{"other", "run"}, goVersionToStderrArgs()[2:]...)
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
		args = []string{"cursor", "run", "cmd", "/C", "echo __FILE__ 1>&2"}
	} else {
		args = []string{"cursor", "run", "sh", "-c", "echo __FILE__ >&2"}
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

func TestExecute_UnknownSubcommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"cursor", "nope"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("expected unknown command message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "wat <host>") {
		t.Fatalf("expected root help after unknown command, got %q", stderr.String())
	}
}

func TestExecute_ParseErrorMissingFilePatternValue(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"cursor", "run", "--file-pattern"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "needs an argument") {
		t.Fatalf("expected flag needs-an-argument for file-pattern, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "__CONVERSATION_ID__") {
		t.Fatalf("expected run help after parse error, got %q", stderr.String())
	}
}

func TestExecute_ParseErrorEmptyFilePatternEquals(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"cursor", "run", "--file-pattern="}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "file-pattern value cannot be empty") {
		t.Fatalf("expected empty pattern message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "__CONVERSATION_ID__") {
		t.Fatalf("expected run help after parse error, got %q", stderr.String())
	}
}

func TestExecute_InvalidFilePatternRegexp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute([]string{"cursor", "run", "-f", `(`, "echo", "x"}, strings.NewReader(validCursorHookJSON()), &stdout, &stderr)
	if code != cli.ExitBadInput {
		t.Fatalf("expected cli.ExitBadInput, got %d", code)
	}
	if !strings.Contains(stderr.String(), "invalid --file-pattern regexp") {
		t.Fatalf("expected regexp error, got %q", stderr.String())
	}
}

func TestExecute_RunWithFilePatternStillRunsSubprocess(t *testing.T) {
	base := goVersionToStderrArgs()
	args := append([]string{base[0], base[1], "-f", `[.]go$`}, base[2:]...)
	stdin := strings.NewReader(validCursorHookJSON())
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Execute(args, stdin, &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected cli.ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
	if !strings.Contains(stderr.String(), "go version") {
		t.Fatalf("expected go version on stderr, got %q", stderr.String())
	}
}

func TestExecute_RunWithFilePatternSkipsWhenPathNoMatch(t *testing.T) {
	jsonTxt := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "c1",
		"generation_id": "g1",
		"model": "m",
		"cursor_version": "1.0",
		"workspace_roots": [],
		"user_email": "",
		"transcript_path": "",
		"file_path": "/repo/file.txt",
		"edits": [{"old_string": "a", "new_string": "b"}]
	}`
	if runtime.GOOS == "windows" {
		var stdout, stderr bytes.Buffer
		code := Execute([]string{"cursor", "run", "-f", `[.]go$`, "cmd", "/C", "echo", "ran", "1>&2"}, strings.NewReader(jsonTxt), &stdout, &stderr)
		if code != cli.ExitSuccess {
			t.Fatalf("expected ExitSuccess, got %d, stderr=%q", code, stderr.String())
		}
		assertHookStdoutJSON(t, stdout.String())
		if strings.Contains(stderr.String(), "ran") {
			t.Fatalf("subprocess should be skipped, stderr=%q", stderr.String())
		}
		return
	}
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"cursor", "run", "-f", `[.]go$`, "sh", "-c", "echo ran >&2"}, strings.NewReader(jsonTxt), &stdout, &stderr)
	if code != cli.ExitSuccess {
		t.Fatalf("expected ExitSuccess, got %d, stderr=%q", code, stderr.String())
	}
	assertHookStdoutJSON(t, stdout.String())
	if strings.Contains(stderr.String(), "ran") {
		t.Fatalf("subprocess should be skipped, stderr=%q", stderr.String())
	}
}
