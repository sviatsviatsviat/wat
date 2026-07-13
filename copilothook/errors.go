package copilothook

import "errors"

// Decode error sentinels for stable error handling across packages.
var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = errors.New("copilothook: empty payload")
	// ErrUnrecognizedFormat indicates the payload is not camelCase or VS Code format.
	ErrUnrecognizedFormat = errors.New("copilothook: decode payload: unrecognized format")
	// ErrEventNameRequired indicates camelCase payloads need WithEvent.
	ErrEventNameRequired = errors.New("copilothook: decode: event name required (camelCase payloads need WithEvent)")
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = errors.New("copilothook: decode payload")
)
