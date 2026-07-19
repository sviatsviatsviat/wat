package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent = model.PostToolFailureEvent

// PostToolFailureHook is the handler context for portable PostToolFailure events.
type PostToolFailureHook = model.PostToolFailureHook

// PostToolFailureResult is the portable hook response for PostToolFailure events.
type PostToolFailureResult = model.PostToolFailureResult

// PostToolFailureResults is the hook-scoped response builder supplied to OnPostToolFailure handlers by registration.
type PostToolFailureResults = model.PostToolFailureResults

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler = model.PostToolFailureHandler

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func (c *chain) OnPostToolFailure(fn PostToolFailureHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterPostToolFailure(c.reg, fn)
	copilot.RegisterPostToolFailure(c.reg, fn)
	cursor.RegisterPostToolFailure(c.reg, fn)
	return c
}
