package cursor_test

import (
	"context"
	"strings"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func ExampleUseHooks() {
	hooks := cursor.UseHooks().BeforeShellExecution(func(ctx context.Context, hook cursor.BeforeShellExecution, r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		if hook.Command == "rm -rf /" {
			return r.Deny("blocked"), nil
		}
		return r.Noop(), nil
	})
	// Hook entrypoint: run.Serve(hooks)
	_ = hooks
	_ = run.Serve
}

// ExampleBeforeMCPExecution gates a remote MCP server by URL. For
// security-critical MCP hooks, set failClosed: true in .cursor/hooks.json so
// crash/timeout/invalid JSON blocks the tool. Cursor defers beforeMCPExecution
// for cloud agents.
func ExampleBeforeMCPExecution() {
	hooks := cursor.UseHooks().BeforeMCPExecution(func(ctx context.Context, hook cursor.BeforeMCPExecution, r cursor.PermissionResults) (cursor.PermissionOutput, error) {
		if hook.URL != "" && !strings.HasPrefix(hook.URL, "https://") {
			return r.Deny("only https MCP servers are allowed"), nil
		}
		if hook.Command != "" && strings.Contains(hook.Command, "untrusted") {
			return r.Deny("stdio MCP command blocked"), nil
		}
		return r.Allow(), nil
	})
	_ = hooks
}
