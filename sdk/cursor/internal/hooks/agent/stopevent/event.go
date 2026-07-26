package stopevent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the stop hook event.
//
// Status is completed, aborted, or error. LoopCount is how many times this stop
// hook has already triggered an automatic follow-up for the conversation
// (starts at 0). Cursor caps auto follow-ups per script with the hooks.json
// loop_limit option (default 5; null means unlimited); that option is not an
// SDK field. The same loop_limit applies to subagentStop follow-ups. Authors
// that return FollowUp should consult LoopCount against that budget.
type Event struct {
	event.Envelope
	// Status is the stop status: completed, aborted, or error.
	Status string `json:"status"`
	// LoopCount is how many automatic follow-ups this stop hook has already
	// triggered for the conversation (starts at 0). See the Event comment for
	// the hooks.json loop_limit interaction.
	LoopCount int `json:"loop_count"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.Stop }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.Stop, hookkit.EventDecoder[Event](c))
}
