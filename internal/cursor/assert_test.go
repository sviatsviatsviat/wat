package cursor

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// cursorHookStdoutSuccessLine is the hook stdout line successful Cursor hooks must write.
// Tests use this value instead of the production DefaultHookResponseLine constant so
// assertions stay explicit if that symbol is renamed or the contract is revisited.
const cursorHookStdoutSuccessLine = "{}\n"

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
