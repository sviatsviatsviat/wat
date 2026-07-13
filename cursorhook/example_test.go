package cursorhook_test

import (
	"context"

	"github.com/sviatsviatsviat/wat/cursorhook"
)

func ExampleMux() {
	mux := cursorhook.NewMux()
	cursorhook.On(mux, func(ctx context.Context, ev cursorhook.BeforeShellExecution) (cursorhook.PermissionOutput, error) {
		if ev.Command == "git push --force" {
			return cursorhook.PermissionOutput{
				Decision:     cursorhook.DecisionDeny,
				AgentMessage: "force push blocked",
			}, nil
		}
		return cursorhook.PermissionOutput{}, nil
	})
	// Hook entrypoint: mux.Main()
}
