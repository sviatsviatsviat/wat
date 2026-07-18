package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// UserPromptExpansion is the UserPromptExpansion hook event.
type UserPromptExpansion struct {
	Envelope
	// ExpansionType is the expansion kind (slash_command, mcp_prompt).
	ExpansionType string `json:"expansion_type"`
	// CommandName is the slash command name.
	CommandName string `json:"command_name"`
	// CommandArgs is the slash command arguments.
	CommandArgs string `json:"command_args"`
	// CommandSource is the command source.
	CommandSource string `json:"command_source"`
	// Prompt is the expanded prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (UserPromptExpansion) EventName() string { return EventUserPromptExpansion }

func init() {
	registerDecoder(EventUserPromptExpansion, decodeAs[UserPromptExpansion])
}

// UserPromptExpansionResults is the hook-scoped response builder supplied to On* handlers by registration.
type UserPromptExpansionResults interface {
	// Context returns a context-injection-only UserPromptExpansion result.
	Context(text string) CommonOutput
	isUserPromptExpansionResults()
}

type userPromptExpansionResults struct{}

func (userPromptExpansionResults) isUserPromptExpansionResults() {}

// Context returns a context-injection-only UserPromptExpansion result.
func (userPromptExpansionResults) Context(text string) CommonOutput {
	return commonOutput{additionalContext: text}
}

// OnUserPromptExpansion registers a UserPromptExpansion handler.
func OnUserPromptExpansion(fn func(context.Context, Hook[UserPromptExpansion], UserPromptExpansionResults) (CommonOutput, error)) *chain {
	return (&chain{}).UserPromptExpansion(fn)
}

// UserPromptExpansion registers another UserPromptExpansion handler on the chain.
func (c *chain) UserPromptExpansion(fn func(context.Context, Hook[UserPromptExpansion], UserPromptExpansionResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev UserPromptExpansion) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), userPromptExpansionResults{})
	})
	return c
}
