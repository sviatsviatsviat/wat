// Package copilot is the GitHub Copilot hook SDK. Hook authors register
// typed handlers; the SDK decodes stdin JSON (camelCase and VS Code compatible
// formats), dispatches by event, encodes stdout, and selects exit codes.
//
// Encode returns flat camelCase stdout JSON and a process exit code. camelCase
// payloads require WithEvent unless hook_event_name is present on the wire.
//
// Example:
//
//	mux := copilot.NewMux()
//	copilot.On(mux, func(ctx context.Context, ev copilot.PreToolUse) (copilot.PreToolOutput, error) {
//	    if ev.NativeToolName() == "bash" {
//	        return copilot.PreToolOutput{Decision: copilot.DecisionDeny, Reason: "blocked"}, nil
//	    }
//	    return copilot.PreToolOutput{}, nil
//	})
//	mux.Main(copilot.WithEvent("preToolUse"))
//
// See ExampleMux in example_test.go for a runnable example.
package copilot
