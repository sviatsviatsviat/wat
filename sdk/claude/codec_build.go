package claude

// Dialect codec helpers. Hook authors use Chain-injected *Results, not these.

// BuildPreToolUseOutput constructs a PreToolUseOutput for dialect codecs.
func BuildPreToolUseOutput(decision PermissionDecision, reason string, updatedInput map[string]any) PreToolUseOutput {
	return preToolUseOutput{decision: decision, reason: reason, updatedInput: updatedInput}
}

// BuildPostToolUseOutput constructs a PostToolUseOutput for dialect codecs.
func BuildPostToolUseOutput(additionalContext string, updatedToolOutput any) PostToolUseOutput {
	return postToolUseOutput{additionalContext: additionalContext, updatedToolOutput: updatedToolOutput}
}

// BuildStopOutput constructs a StopOutput for dialect codecs.
func BuildStopOutput(followUp string) StopOutput {
	if followUp == "" {
		return stopOutput{}
	}
	return stopOutput{block: true, reason: followUp}
}

// BuildSessionStartOutput constructs a SessionStartOutput for dialect codecs.
func BuildSessionStartOutput(additionalContext string) SessionStartOutput {
	return sessionStartOutput{additionalContext: additionalContext}
}

// BuildPostToolUseFailureOutput constructs a PostToolUseFailure context output for dialect codecs.
func BuildPostToolUseFailureOutput(additionalContext string) PostToolUseOutput {
	return postToolUseOutput{additionalContext: additionalContext}
}
