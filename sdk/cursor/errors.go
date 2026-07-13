package cursor

import "errors"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("cursor: empty payload")
	// ErrEventNameRequired indicates payloads need hook_event_name or WithEvent.
	ErrEventNameRequired = errors.New("cursor: decode: event name required (use WithEvent when hook_event_name is absent)")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("cursor: decode payload")
)
