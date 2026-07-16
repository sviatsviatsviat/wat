package copilot

// Dialect codec helpers. Hook authors use Chain-injected *Results, not these.

// BuildPreToolOutput constructs a PreToolOutput for dialect codecs.
func BuildPreToolOutput(decision PermissionDecision, reason string, modifiedArgs map[string]any) PreToolOutput {
	return preToolOutput{decision: decision, reason: reason, modifiedArgs: modifiedArgs}
}

// BuildPostToolOutput constructs a PostToolOutput for dialect codecs.
func BuildPostToolOutput(additionalContext, modifiedResult string) PostToolOutput {
	return postToolOutput{additionalContext: additionalContext, modifiedResult: modifiedResult}
}

// BuildStopOutput constructs a StopOutput for dialect codecs.
func BuildStopOutput(reason string) StopOutput {
	return stopOutput{reason: reason}
}

// BuildSessionStartOutput constructs a SessionStartOutput for dialect codecs.
func BuildSessionStartOutput(additionalContext string) SessionStartOutput {
	return sessionStartOutput{additionalContext: additionalContext}
}

// BuildPostToolFailureOutput constructs a PostToolFailureOutput for dialect codecs.
func BuildPostToolFailureOutput(context string) PostToolFailureOutput {
	return postToolFailureOutput{context: context}
}
