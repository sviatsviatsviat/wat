package erroroccurred

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Event is the ErrorOccurred hook event.
type Event struct {
	event.Envelope
	// Error is the error payload (string or object).
	Error json.RawMessage `json:"error"`
	// ErrorContext is additional error context.
	ErrorContext string `json:"error_context"`
	// Recoverable is true when the error may be retried.
	Recoverable *bool `json:"recoverable"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.ErrorOccurred }

// Context returns error context.
func (e Event) Context() string {
	return e.ErrorContext
}

// Detail parses structured error details when present.
func (e Event) Detail() (event.ErrorDetail, bool) {
	if hookkit.NullToNil(e.Error) == nil {
		return event.ErrorDetail{}, false
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return event.ErrorDetail{Message: s}, true
	}
	var detail event.ErrorDetail
	if json.Unmarshal(e.Error, &detail) != nil {
		return event.ErrorDetail{}, false
	}
	return detail, true
}

// Register registers this hook event decoder on c.
func Register(c *hookkit.Codec) {
	c.Register(event.ErrorOccurred, hookkit.EventDecoder[Event](c))
}

// RegisterHandler registers an ErrorOccurred observe handler on reg.
func RegisterHandler(reg *run.Registry, fn func(context.Context, run.Hook[Event]) error) {
	if fn == nil {
		return
	}
	reg.RegisterObserveHandler(runtime.Dialect, run.ObserveHandler(fn))
}
