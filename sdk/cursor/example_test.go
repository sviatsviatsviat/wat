package cursor_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleUseHooks() {
	cursor.UseHooks().BeforeShellExecution(func(ctx context.Context, hook cursor.BeforeShellExecution, r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		if hook.Command == "rm -rf /" {
			return r.Deny("blocked"), nil
		}
		return r.Noop(), nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
