package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PostToolFailureEvent is the normalized view of a PostToolFailure hook invocation.
type PostToolFailureEvent = model.PostToolFailureEvent

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
	claude.RegisterPostToolFailure(fn)
	copilot.RegisterPostToolFailure(fn)
	cursor.RegisterPostToolFailure(fn)
	return c
}
