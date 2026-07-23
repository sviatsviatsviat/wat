package runtime

import (
	"errors"
)

// Dialect is the sdk/run router name for Claude Code hooks.
const Dialect = "claude"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("claude: empty payload")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("claude: decode payload")
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = errors.New("claude: decode: event name required (hook_event_name is required)")
)
