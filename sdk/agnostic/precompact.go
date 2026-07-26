package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
// Portable CompactInfo exposes only fields all three agents share (trigger and,
// when present, custom instructions). Cursor-only compaction metrics and
// user_message output require sdk/cursor.PreCompact.
type PreCompactEvent = model.PreCompactEvent

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler = model.PreCompactHandler

// OnPreCompact registers an observe-only handler for PreCompact events.
// Unlike sdk/cursor.PreCompact, this portable registration cannot emit a
// user_message; use the native Cursor registrar for that capability.
func (c *hooks) OnPreCompact(fn PreCompactHandler) *hooks {
	if fn == nil {
		return c
	}
	return c.appendParts(
		claude.RegisterPreCompact(fn),
		copilot.RegisterPreCompact(fn),
		cursor.RegisterPreCompact(fn),
	)
}
