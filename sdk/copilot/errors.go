package copilot

import "errors"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("copilot: empty payload")
	// ErrUnrecognizedFormat indicates the payload is not camelCase or VS Code format.
	ErrUnrecognizedFormat = errors.New("copilot: decode payload: unrecognized format")
	// ErrEventNameRequired indicates camelCase payloads need WithEvent.
	ErrEventNameRequired = errors.New("copilot: decode: event name required (camelCase payloads need WithEvent)")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("copilot: decode payload")
)
