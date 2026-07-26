package afteragentthought

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the afterAgentThought hook event. It is observe-only: Cursor does
// not currently support output fields for this hook.
//
// Cursor's hooks.json matcher for afterAgentThought is the fixed value
// AgentThought. That filter is host configuration, not a field on this
// payload.
type Event struct {
	event.Envelope
	// Text is the fully aggregated thinking text for the completed block.
	Text string `json:"text"`
	// DurationMs is the thinking-block duration in milliseconds when present.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.AfterAgentThought }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.AfterAgentThought, hookkit.EventDecoder[Event](c))
}
