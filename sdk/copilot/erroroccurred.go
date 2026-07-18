package copilot

import (
	"context"
	"encoding/json"
)

// ErrorOccurred is the errorOccurred hook event.
type ErrorOccurred struct {
	Envelope
	// Error is the error payload (string or object).
	Error json.RawMessage `json:"error"`
	// ErrorContext is additional error context (VS Code).
	ErrorContext string `json:"error_context"`
	// ErrorContextCamel is additional error context (camelCase).
	ErrorContextCamel string `json:"errorContext"`
	// Recoverable is true when the error may be retried.
	Recoverable *bool `json:"recoverable"`
}

// EventName returns the canonical hook event name.
func (ErrorOccurred) EventName() string { return EventErrorOccurred }

// Context returns error context from either wire format.
func (e ErrorOccurred) Context() string {
	if e.ErrorContext != "" {
		return e.ErrorContext
	}
	return e.ErrorContextCamel
}

// Detail parses structured error details when present.
func (e ErrorOccurred) Detail() (ErrorDetail, bool) {
	if len(e.Error) == 0 || string(bytesTrimSpace(e.Error)) == "null" {
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
	registerDecoder(EventErrorOccurred, decodeAs[ErrorOccurred])
}

// OnErrorOccurred registers an observe-only errorOccurred handler.
func OnErrorOccurred(fn func(context.Context, Hook[ErrorOccurred]) error) *chain {
	return (&chain{}).ErrorOccurred(fn)
}

// ErrorOccurred registers another ErrorOccurred handler on the chain.
func (c *chain) ErrorOccurred(fn func(context.Context, Hook[ErrorOccurred]) error) *chain {
	registerObserveHandler(fn)
	return c
}
