package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func assertEqual(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func assertIntEqual(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Fatalf("want %d, got %d", want, got)
	}
}

func assertTemplateBindingValue(t *testing.T, bindings interface {
	TemplateValue(key string) (string, bool)
}, key, want string,
) {
	t.Helper()
	bindingValue, ok := bindings.TemplateValue(key)
	if !ok {
		t.Fatalf("TemplateValue(%q): expected ok true", key)
	}
	if bindingValue != want {
		t.Fatalf("TemplateValue(%q): want %q, got %q", key, want, bindingValue)
	}
}

// stubHookCommand implements [core.Command] for hook handler tests.
type stubHookCommand struct {
	execute func(*core.HookContext) int
}

func (stub stubHookCommand) Execute(ctx *core.HookContext) int {
	if stub.execute == nil {
		return 0
	}
	return stub.execute(ctx)
}
