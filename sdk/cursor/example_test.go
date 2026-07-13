package cursor_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

func ExampleMux() {
	mux := cursor.NewMux()
	cursor.On(mux, func(ctx context.Context, ev cursor.BeforeShellExecution) (cursor.PermissionOutput, error) {
		if ev.Command == "git push --force" {
			return cursor.PermissionOutput{
				Decision:     cursor.DecisionDeny,
				AgentMessage: "force push blocked",
			}, nil
		}
		return cursor.PermissionOutput{}, nil
	})
	// Hook entrypoint: mux.Main()
}
