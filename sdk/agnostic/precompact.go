package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent = model.PreCompactEvent

// PreCompactHook is the handler context for portable PreCompact events.
type PreCompactHook = model.PreCompactHook

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler = model.PreCompactHandler

// OnPreCompact registers an observe-only handler for PreCompact events.
func OnPreCompact(fn PreCompactHandler) *chain {
	if fn == nil {
		return &chain{}
	}
	claude.RegisterPreCompact(fn)
	copilot.RegisterPreCompact(fn)
	cursor.RegisterPreCompact(fn)
	return &chain{}
}

// OnPreCompact registers another observe-only PreCompact handler on the chain.
func (c *chain) OnPreCompact(fn PreCompactHandler) *chain {
	return OnPreCompact(fn)
}
