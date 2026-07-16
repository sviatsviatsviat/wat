package cursor

// Dialect codec helpers. Hook authors use Chain-injected *Results, not these.

// BuildPermissionOutput constructs a PermissionOutput for dialect codecs.
func BuildPermissionOutput(decision PermissionDecision, agentMessage string, updatedInput map[string]any) PermissionOutput {
	return permissionOutput{decision: decision, agentMessage: agentMessage, updatedInput: updatedInput}
}

// BuildPostToolOutput constructs a PostToolOutput for dialect codecs.
func BuildPostToolOutput(additionalContext string, updatedMCPOutput any) PostToolOutput {
	return postToolOutput{additionalContext: additionalContext, updatedMCPOutput: updatedMCPOutput}
}

// BuildStopOutput constructs a StopOutput for dialect codecs.
func BuildStopOutput(followUpMessage string) StopOutput {
	return stopOutput{followUpMessage: followUpMessage}
}

// BuildSessionStartOutput constructs a SessionStartOutput for dialect codecs.
func BuildSessionStartOutput(additionalContext string) SessionStartOutput {
	return sessionStartOutput{additionalContext: additionalContext}
}

// BuildBeforeSubmitPromptOutput constructs a BeforeSubmitPromptOutput for dialect codecs.
func BuildBeforeSubmitPromptOutput(cont *bool, userMessage string) BeforeSubmitPromptOutput {
	return beforeSubmitPromptOutput{cont: cont, userMessage: userMessage}
}

// BuildPreCompactOutput constructs a PreCompactOutput for dialect codecs.
func BuildPreCompactOutput(userMessage string) PreCompactOutput {
	return preCompactOutput{userMessage: userMessage}
}
