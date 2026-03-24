package app

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestNewHookHandlerFactory_Cursor(t *testing.T) {
	execCtx := core.NewWatExecutionContext("cursor").WithSubcommand("run")
	factory, err := newHookHandlerFactory(execCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
}

func TestNewHookHandlerFactory_UnsupportedHost(t *testing.T) {
	execCtx := core.NewWatExecutionContext("other")
	factory, err := newHookHandlerFactory(execCtx)
	if err == nil {
		t.Fatal("expected error")
	}
	if factory != nil {
		t.Fatalf("expected nil factory, got %#v", factory)
	}
	if !strings.Contains(err.Error(), `host "other" is not supported yet`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
