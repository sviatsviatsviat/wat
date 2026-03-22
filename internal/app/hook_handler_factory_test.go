package app

import (
	"strings"
	"testing"
)

func TestNewHookHandlerFactory_Cursor(t *testing.T) {
	factory, err := newHookHandlerFactory("cursor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if factory == nil {
		t.Fatal("expected non-nil factory")
	}
}

func TestNewHookHandlerFactory_UnsupportedHost(t *testing.T) {
	factory, err := newHookHandlerFactory("other")
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
