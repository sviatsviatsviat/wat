// Package claudehook is the Claude Code hook SDK. Hook authors register typed
// handlers with On and call Main; the SDK decodes stdin JSON, dispatches by
// event, encodes stdout, and selects exit codes per the Claude Code hook protocol.
//
// Example:
//
//	mux := claudehook.NewMux()
//	claudehook.On(mux, func(ctx context.Context, ev claudehook.PreToolUse) (claudehook.PreToolUseOutput, error) {
//	    if ev.ToolName == "Bash" {
//	        return claudehook.PreToolUseOutput{Decision: claudehook.DecisionDeny, Reason: "blocked"}, nil
//	    }
//	    return claudehook.PreToolUseOutput{}, nil
//	})
//	mux.Main()
package claudehook
