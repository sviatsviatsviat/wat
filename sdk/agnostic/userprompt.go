package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent = model.UserPromptEvent

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler = model.UserPromptHandler

// OnUserPrompt registers an observe-only handler for UserPrompt events.
func (c *hooks) OnUserPrompt(fn UserPromptHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterUserPrompt(fn),
		copilot.RegisterUserPrompt(fn),
		cursor.RegisterUserPrompt(fn),
	)
}
