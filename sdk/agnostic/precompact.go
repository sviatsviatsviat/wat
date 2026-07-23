package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent = model.PreCompactEvent

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler = model.PreCompactHandler

// OnPreCompact registers an observe-only handler for PreCompact events.
func (c *chain) OnPreCompact(fn PreCompactHandler) *chain {
	if fn == nil {
		return c
	}
	claude.RegisterPreCompact(fn)
	copilot.RegisterPreCompact(fn)
	cursor.RegisterPreCompact(fn)
	return c
}
