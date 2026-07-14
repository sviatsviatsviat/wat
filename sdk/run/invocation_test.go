package run

import (
	"context"
	"testing"
)

func TestInvocationFrom(t *testing.T) {
	ctx := WithConfig(context.Background(), &Config{
		Dialect:   "claude",
		EventHint: "PreToolUse",
		Getenv:    func(key string) string { return "v:" + key },
	})
	inv := InvocationFrom(ctx)
	if inv.Dialect() != "claude" {
		t.Fatalf("Dialect() = %q", inv.Dialect())
	}
	if inv.EventHint() != "PreToolUse" {
		t.Fatalf("EventHint() = %q", inv.EventHint())
	}
	if got := inv.Getenv("WAT_AGENT"); got != "v:WAT_AGENT" {
		t.Fatalf("Getenv() = %q", got)
	}
}

func TestInvocationFrom_defaultContext(t *testing.T) {
	inv := InvocationFrom(context.Background())
	if inv.Dialect() != "" || inv.EventHint() != "" {
		t.Fatalf("Dialect=%q EventHint=%q", inv.Dialect(), inv.EventHint())
	}
}
