package cursor_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleOn() {
	cursor.On(func(ctx context.Context, ev cursor.BeforeShellExecution) (cursor.PermissionOutput, error) {
		if ev.Command == "rm -rf /" {
			return cursor.PermissionOutput{Decision: cursor.DecisionDeny, AgentMessage: "blocked"}, nil
		}
		return cursor.PermissionOutput{}, nil
	})
	// Hook entrypoint: run.Main()
	_ = run.Main
}
