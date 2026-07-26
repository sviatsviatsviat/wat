package sessionend

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the sessionEnd hook event.
//
// sessionEnd is fire-and-forget: handlers observe only and the host ignores
// any response body. Cursor documents this hook for IDE composer sessions; it
// is not available for cloud agents.
type Event struct {
	event.Envelope
	// Reason is how the session ended (for example completed, aborted, error,
	// window_close, or user_close).
	Reason string `json:"reason"`
	// DurationMs is the total session duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	// IsBackgroundAgent reports whether this was a background agent session.
	IsBackgroundAgent bool `json:"is_background_agent"`
	// FinalStatus is the final status of the session.
	FinalStatus string `json:"final_status"`
	// ErrorMessage is the error details when Reason is "error".
	ErrorMessage string `json:"error_message"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.SessionEnd }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.SessionEnd, hookkit.EventDecoder[Event](c))
}
