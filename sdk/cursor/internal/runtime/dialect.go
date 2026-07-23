package runtime

import (
	"errors"
)

// Dialect is the sdk/run registry name for Cursor hooks.
const Dialect = "cursor"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("cursor: empty payload")
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = errors.New("cursor: decode: event name required (hook_event_name is required)")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("cursor: decode payload")
)
