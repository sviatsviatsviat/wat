package runtime

import (
	"errors"
)

// Dialect is the sdk/run registry name for GitHub Copilot hooks.
const Dialect = "copilot"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("copilot: empty payload")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("copilot: decode payload")
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = errors.New("copilot: decode: event name required (hook_event_name is required)")
)
