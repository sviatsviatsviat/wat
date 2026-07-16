package cursor_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleChain() {
	new(cursor.Chain).BeforeShellExecution(func(ctx context.Context, hook cursor.Hook[cursor.BeforeShellExecution], r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		if hook.Event.Command == "rm -rf /" {
			return r.Deny("blocked"), nil
		}
		return nil, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
