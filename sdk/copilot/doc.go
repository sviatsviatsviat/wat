// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers with UseHooks into a run.Registry (the process default when
// called with no arguments), then call run.Main from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Hook stdout JSON is encoded internally from Output (sealed; only this package
// implements it) as snake_case JSON with a process exit code. Payloads must
// include hook_event_name on the wire.
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
