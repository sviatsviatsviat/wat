package execcommand

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func ptr(s string) *string { return &s }

func TestTemplateBindingsAfterFileEdit_definedKeys(t *testing.T) {
	h := testAfterFileEditHookStub{
		hookEventName:  "afterFileEdit",
		transcriptPath: ptr("/t.jsonl"),
		filePath:       "D:/repo/x.go",
	}
	b := templateBindingsAfterFileEdit{hook: &h}
	assertTemplateBindingValue(t, b, "HOOK_EVENT_NAME", "afterFileEdit")
	assertTemplateBindingValue(t, b, "TRANSCRIPT_PATH", "/t.jsonl")
	assertTemplateBindingValue(t, b, "FILE_PATH", "D:/repo/x.go")
	_, ok := b.TemplateValue("DURATION")
	if ok {
		t.Fatal("DURATION must not be defined for after file edit bindings")
	}
}

func TestTemplateBindingsAfterFileEdit_nullTranscript(t *testing.T) {
	h := testAfterFileEditHookStub{
		hookEventName: "afterTabFileEdit",
		filePath:      "/x",
	}
	b := templateBindingsAfterFileEdit{hook: &h}
	assertTemplateBindingValue(t, b, "TRANSCRIPT_PATH", "")
}

func TestTemplateBindingsAfterShellExecution_definedKeys(t *testing.T) {
	h := testAfterShellHookStub{
		hookEventName:  "afterShellExecution",
		transcriptPath: ptr("/t.jsonl"),
		duration:       12.5,
		rawCommand:     "go test ./...",
		sandbox:        true,
	}
	b := templateBindingsAfterShellExecution{hook: &h}
	assertTemplateBindingValue(t, b, "HOOK_EVENT_NAME", "afterShellExecution")
	assertTemplateBindingValue(t, b, "TRANSCRIPT_PATH", "/t.jsonl")
	assertTemplateBindingValue(t, b, "DURATION", "12.5")
	assertTemplateBindingValue(t, b, "SANDBOX", "true")
	assertTemplateBindingValue(t, b, "COMMAND", "go test ./...")
	_, ok := b.TemplateValue("FILE_PATH")
	if ok {
		t.Fatal("FILE_PATH must not be defined for shell bindings")
	}
}

func TestTemplateBindingsAfterShellExecution_decimalDuration(t *testing.T) {
	h := testAfterShellHookStub{
		hookEventName: "afterShellExecution",
		duration:      2841.805,
	}
	b := templateBindingsAfterShellExecution{hook: &h}
	assertTemplateBindingValue(t, b, "DURATION", "2841.805")
}

type testAfterFileEditHookStub struct {
	hookEventName  string
	transcriptPath *string
	filePath       string
}

func (testAfterFileEditHookStub) WriteDefaultToHost() {}

func (s *testAfterFileEditHookStub) HookEventName() string { return s.hookEventName }

func (s *testAfterFileEditHookStub) TranscriptPath() *string { return s.transcriptPath }

func (s *testAfterFileEditHookStub) FilePath() string { return s.filePath }

var _ core.AfterFileEditHook = (*testAfterFileEditHookStub)(nil)

type testAfterShellHookStub struct {
	hookEventName  string
	transcriptPath *string
	duration       float32
	rawCommand     string
	sandbox        bool
}

func (testAfterShellHookStub) WriteDefaultToHost() {}

func (s *testAfterShellHookStub) HookEventName() string { return s.hookEventName }

func (s *testAfterShellHookStub) TranscriptPath() *string { return s.transcriptPath }

func (s *testAfterShellHookStub) Duration() float32 { return s.duration }

func (s *testAfterShellHookStub) RawCommand() string { return s.rawCommand }

func (s *testAfterShellHookStub) Sandbox() bool { return s.sandbox }

var _ core.AfterShellExecutionHook = (*testAfterShellHookStub)(nil)

func assertTemplateBindingValue(t *testing.T, bindings templateBindings, key, want string) {
	t.Helper()
	bindingValue, ok := bindings.TemplateValue(key)
	if !ok {
		t.Fatalf("TemplateValue(%q): expected ok true", key)
	}
	if bindingValue != want {
		t.Fatalf("TemplateValue(%q): want %q, got %q", key, want, bindingValue)
	}
}
