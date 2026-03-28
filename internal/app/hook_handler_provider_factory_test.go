package app

import (
	"errors"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cli"
)

func TestNewHookHandlerProvider_Exec(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := newHookHandlerProvider("exec", mockConsole, []string{"echo", "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestNewHookHandlerProvider_ExecWithFilePatternFlag(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := newHookHandlerProvider("exec", mockConsole, []string{"-f", `[.]go$`, "echo", "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestNewHookHandlerProvider_ExecEmptyArgs(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := newHookHandlerProvider("exec", mockConsole, []string{})
	if err == nil {
		t.Fatal("expected error from execcommand.NewExecHookHandlerProvider")
	}
	if provider != nil {
		t.Fatal("expected nil provider")
	}
	if errors.Is(err, errHookHandlerProviderBadInput) {
		t.Fatal("exec parse error should not be errHookHandlerProviderBadInput")
	}
}

func TestNewHookHandlerProvider_UnknownSubcommand(t *testing.T) {
	mockConsole := cli.NewMockConsole()
	provider, err := newHookHandlerProvider("nope", mockConsole, []string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if provider != nil {
		t.Fatalf("expected nil provider, got %#v", provider)
	}
	if !errors.Is(err, errHookHandlerProviderBadInput) {
		t.Fatalf("expected errHookHandlerProviderBadInput, got %v", err)
	}
}
