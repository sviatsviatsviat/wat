package sessionstart

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the sessionStart hook event.
//
// Cursor runs sessionStart as fire-and-forget: the agent loop does not wait
// for or enforce a blocking response. The wire schema also accepts continue
// and user_message output fields, but current callers do not enforce them;
// session creation proceeds even when continue is false. Prefer env and
// additional_context via [Results] and [Output].
//
// sessionStart is not available for Cursor cloud agents.
type Event struct {
	event.Envelope
	// IsBackgroundAgent reports whether this is a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
	// ComposerMode is the composer mode when present on the wire
	// ("agent", "ask", or "edit").
	ComposerMode string `json:"composer_mode"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionStart }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SessionStart, hookkit.EventDecoder[Event](c))
}
