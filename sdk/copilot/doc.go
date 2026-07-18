// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers with On* helpers into the shared run registry, then call
// run.Main from github.com/sviatsviatsviat/wat/sdk/run with run.WithEvent when
// needed.
//
// Encode returns flat camelCase stdout JSON and a process exit code. camelCase
// payloads require run.WithEvent unless hook_event_name is present on the wire.
//
// Example:
//
//	copilot.OnPreToolUse(func(ctx context.Context, hook copilot.Hook[copilot.PreToolUse], r copilot.PreToolResults) (copilot.PreToolOutput, error) {
//	    if hook.Event.NativeToolName() == "bash" {
//	        return r.Deny("blocked"), nil
//	    }
//	    return r.Noop(), nil
//	})
//	run.Main(run.WithEvent("preToolUse"))
//
// See ExampleOnPreToolUse in example_test.go for a runnable example.
package copilot
