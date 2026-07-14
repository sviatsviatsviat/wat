// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers with On into the shared run registry, then call run.Main
// from github.com/sviatsviatsviat/wat/sdk/run with run.WithEvent when needed.
//
// Encode returns flat camelCase stdout JSON and a process exit code. camelCase
// payloads require run.WithEvent (or Decode WithEvent) unless hook_event_name
// is present on the wire.
//
// Example:
//
//	copilot.On(func(ctx context.Context, ev copilot.PreToolUse) (copilot.PreToolOutput, error) {
//	    if ev.NativeToolName() == "bash" {
//	        return copilot.PreToolOutput{Decision: copilot.DecisionDeny, Reason: "blocked"}, nil
//	    }
//	    return copilot.PreToolOutput{}, nil
//	})
//	run.Main(run.WithEvent("preToolUse"))
//
// See ExampleOn in example_test.go for a runnable example.
package copilot
