package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent = model.UserPromptEvent

// UserPromptHook is the handler context for portable UserPrompt events.
type UserPromptHook = model.UserPromptHook

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler = model.UserPromptHandler

// OnUserPrompt registers an observe-only handler for UserPrompt events.
func OnUserPrompt(fn UserPromptHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterUserPrompt(fn)
	copilot.RegisterUserPrompt(fn)
	cursor.RegisterUserPrompt(fn)
	return &Chain{}
}

// OnUserPrompt registers another observe-only UserPrompt handler on the chain.
func (c *Chain) OnUserPrompt(fn UserPromptHandler) *Chain {
	return OnUserPrompt(fn)
}
