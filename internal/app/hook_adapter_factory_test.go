package app

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestNewHookAdapterFactory_Cursor(t *testing.T) {
	factory, err := newHookAdapterFactory("cursor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
	var _ core.HookAdapterFactory = factory
}

func TestNewHookAdapterFactory_UnsupportedHost(t *testing.T) {
	factory, err := newHookAdapterFactory("other")
	if err == nil {
		t.Fatal("expected error")
	}
	if factory != nil {
		t.Fatal("expected nil factory")
	}
}
