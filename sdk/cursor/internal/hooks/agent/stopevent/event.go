package stopevent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the stop hook event.
//
// Follow-up output is built via Results.FollowUp. Host-side auto follow-up
// caps use the hooks.json loop_limit handler option (default 5); that option
// is not an SDK field. The same loop_limit applies to subagentStop follow-ups.
type Event struct {
	event.Envelope
	// Status is the stop status.
	Status string `json:"status"`
	// LoopCount is how many times the stop hook has already triggered an
	// automatic follow-up for this conversation (starts at 0).
	LoopCount int `json:"loop_count"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.Stop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Stop, hookkit.EventDecoder[Event](c))
}
