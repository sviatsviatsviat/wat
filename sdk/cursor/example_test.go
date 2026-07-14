package cursor_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleChain() {
	new(cursor.Chain).BeforeShellExecution(func(ctx context.Context, hook cursor.BeforeShellExecutionHook, r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		if hook.Event.Command == "rm -rf /" {
			return r.Deny("blocked"), nil
		}
		return cursor.PermissionOutput{}, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
