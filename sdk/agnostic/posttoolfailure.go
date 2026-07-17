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

// PostToolFailureResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PostToolFailureResults = model.PostToolFailureResults

// PostToolFailureHandler handles portable PostToolFailure events.
type PostToolFailureHandler = model.PostToolFailureHandler

// OnPostToolFailure registers a handler for PostToolFailure events across all agents.
func OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	claude.RegisterPostToolFailure(fn)
	copilot.RegisterPostToolFailure(fn)
	cursor.RegisterPostToolFailure(fn)
	return &Chain{}
}

// OnPostToolFailure registers another PostToolFailure handler on the chain.
func (c *Chain) OnPostToolFailure(fn PostToolFailureHandler) *Chain {
	return OnPostToolFailure(fn)
}
