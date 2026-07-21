package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// ErrorOccurred is the ErrorOccurred hook event.
type ErrorOccurred struct {
	Envelope
	// Error is the error payload (string or object).
	Error json.RawMessage `json:"error"`
	// ErrorContext is additional error context.
	ErrorContext string `json:"error_context"`
	// Recoverable is true when the error may be retried.
	Recoverable *bool `json:"recoverable"`
}

// EventName returns the canonical hook event name.
func (ErrorOccurred) EventName() string { return EventErrorOccurred }

// Context returns error context.
func (e ErrorOccurred) Context() string {
	return e.ErrorContext
}

// Detail parses structured error details when present.
func (e ErrorOccurred) Detail() (ErrorDetail, bool) {
	if hookkit.NullToNil(e.Error) == nil {
		return ErrorDetail{}, false
	}
	var s string
	if json.Unmarshal(e.Error, &s) == nil {
		return ErrorDetail{Message: s}, true
	}
	var detail ErrorDetail
	if json.Unmarshal(e.Error, &detail) != nil {
		return ErrorDetail{}, false
	}
	return detail, true
}

func init() {
	codec.Register(EventErrorOccurred, hookkit.EventDecoder[ErrorOccurred](codec))
}

// ErrorOccurred registers a ErrorOccurred handler on the chain.
func (c *chain) ErrorOccurred(fn func(context.Context, run.Hook[ErrorOccurred]) error) *chain {
	c.reg.RegisterObserveHandler(Dialect, run.ObserveHandler(fn))
	return c
}
