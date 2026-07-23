// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers with UseHooks, then call run.Main from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Hook stdout JSON is encoded internally from Output (sealed; only this package
// implements it) as snake_case JSON with a process exit code. Payloads must
// include hook_event_name on the wire.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsBash,
// AsCreate, and related accessors (and Tool* name constants) in this package.
//
// Example:
//
//	copilot.UseHooks().PreToolUse(func(ctx context.Context, hook run.Hook[copilot.PreToolUse], r copilot.PreToolResults) (copilot.PreToolOutput, error) {
//	    if hook.Event.NativeToolName() == "bash" {
//	        return r.Deny("blocked"), nil
//	    }
//	    return r.Noop(), nil
//	})
//	run.Main()
//
// See ExampleUseHooks in example_test.go for a runnable example.
package copilot
