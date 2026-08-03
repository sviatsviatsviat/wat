// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers with UseHooks, then pass the chain to run.Serve from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// The supported host wire is the Copilot hooks-reference PascalCase event
// names with snake_case JSON fields (flat stdout). Copilot CLI camelCase
// samples and VS Code agent-customization hookSpecificOutput are not this
// dialect; see docs/agents/copilot.md.
//
// Hook stdout JSON is encoded internally from Output (sealed; only this package
// implements it) as snake_case JSON with a process exit code. Payloads must
// include hook_event_name on the wire.
//
// PermissionRequest.SoftDeny is a soft deny (behavior "deny", exit 0), not a user
// confirmation prompt: Copilot's permissionRequest schema accepts only allow or
// deny. Prefer Deny to block, or Noop/nil to fall through to the host permission
// UI. Under Copilot cloud agent, permissionRequest is inert (use PreToolUse);
// PreToolUse "ask" is treated as deny; command-hook timeouts are fail-open.
// See docs/agents/copilot.md.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsBash,
// AsCreate, AsEdit, AsView, AsGlob, AsGrep, AsWebFetch, AsWebSearch, AsTask,
// AsAskUser, AsUpdateTodo, and related Tool* name constants in this package.
// Accessors accept runtime names and Claude-format names used on PascalCase
// PreToolUse payloads (for example Bash, Edit, Glob).

// Example:
//
//	run.Serve(
//	    copilot.UseHooks().PreToolUse(func(ctx context.Context, hook copilot.PreToolUse, r copilot.PreToolResults) (copilot.PreToolOutput, error) {
//	        if hook.NativeToolName() == "bash" {
//	            return r.Deny("blocked"), nil
//	        }
//	        return r.Noop(), nil
//	    }),
//	)
//
// See ExampleUseHooks in example_test.go for a runnable example.
package copilot
