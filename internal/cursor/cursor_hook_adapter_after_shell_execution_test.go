package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/cmdast"
)

func TestNewAfterShellExecutionHookAdapter_success(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"npm test","output":"ok","duration":1,"sandbox":false}`)
	factory := NewHookAdapterFactory()
	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	if adapter == nil {
		t.Fatal("expected non-nil HookAdapter")
	}
}

func TestAfterShellExecutionHookAdapter_carriesHookDataAndProtocol(t *testing.T) {
	mock := cli.NewMockConsole()
	raw := []byte(`{"hook_event_name":"afterShellExecution","conversation_id":"cid-1","command":"cd /tmp/example && git status","output":"all good","duration":1234,"sandbox":true}`)
	factory := NewHookAdapterFactory()
	a, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	adapter, ok := a.(*AfterShellExecutionCursorHookAdapter)
	if !ok || adapter == nil {
		t.Fatalf("want *AfterShellExecutionCursorHookAdapter, got %T", a)
	}
	if adapter.CommonInput.HookEventName != "afterShellExecution" {
		t.Errorf("HOOK_EVENT_NAME: got %q", adapter.CommonInput.HookEventName)
	}
	if adapter.CommonInput.ConversationID != "cid-1" {
		t.Errorf("CONVERSATION_ID: got %q", adapter.CommonInput.ConversationID)
	}
	if adapter.EventSpecificInput == nil {
		t.Fatal("EventSpecificInput must be set")
	}
	if adapter.EventSpecificInput.Command != "cd /tmp/example && git status" {
		t.Errorf("COMMAND: got %q", adapter.EventSpecificInput.Command)
	}
	if adapter.EventSpecificInput.CommandLineErr != nil {
		t.Errorf("CommandLineErr: %v", adapter.EventSpecificInput.CommandLineErr)
	}
	cl := adapter.EventSpecificInput.CommandLine
	if cl == nil {
		t.Fatal("CommandLine must be set")
	}
	if cl.Raw != adapter.EventSpecificInput.Command {
		t.Errorf("CommandLine.Raw: got %q", cl.Raw)
	}
	if len(cl.Statements) != 1 {
		t.Fatalf("statements: %d", len(cl.Statements))
	}
	root := cl.Statements[0]
	if root.Kind != cmdast.StmtChain || root.Chain == nil || root.Chain.Operator != "&&" {
		t.Fatalf("want && chain root: %+v", root)
	}
	if root.Chain.Left == nil || root.Chain.Left.Command == nil || root.Chain.Left.Command.Name != "cd" {
		t.Fatalf("left cmd: %+v", root.Chain.Left)
	}
	if root.Chain.Right == nil || root.Chain.Right.Command == nil || root.Chain.Right.Command.Name != "git" {
		t.Fatalf("right cmd: %+v", root.Chain.Right)
	}
	if adapter.EventSpecificInput.Output != "all good" {
		t.Errorf("OUTPUT: got %q", adapter.EventSpecificInput.Output)
	}
	if adapter.EventSpecificInput.Duration != 1234 {
		t.Errorf("DURATION: got %v", adapter.EventSpecificInput.Duration)
	}
	if !adapter.EventSpecificInput.Sandbox {
		t.Error("SANDBOX: want true")
	}
	adapter.ReturnEmpty()
	if mock.StdoutString() != cursorHookStdoutSuccessLine {
		t.Fatalf("hook stdout after ReturnEmpty: want %q, got %q", cursorHookStdoutSuccessLine, mock.StdoutString())
	}
}

func TestHookAdapterFactory_afterShellExecutionUsesCursorHookAdapter(t *testing.T) {
	mock := cli.NewMockConsole()
	factory := NewHookAdapterFactory()
	raw := []byte(`{"hook_event_name":"afterShellExecution","command":"x","output":"","duration":0,"sandbox":false}`)

	adapter, err := factory.HookAdapterFromJSON(raw, mock)
	if err != nil {
		t.Fatalf("HookAdapterFromJSON: %v", err)
	}
	if _, ok := adapter.(*AfterShellExecutionCursorHookAdapter); !ok {
		t.Fatalf("adapter type: want *AfterShellExecutionCursorHookAdapter, got %T", adapter)
	}
}
