package hookkit

import (
	"encoding/json"
	"fmt"
)

// DecodeAsAndThen unmarshals raw into T, then optionally runs after with the
// original stdin bytes (for tool-input construction and similar decode-time work).
func DecodeAsAndThen[T any](raw []byte, after func(*T, []byte)) (T, error) {
	var ev T
	if err := json.Unmarshal(raw, &ev); err != nil {
		return ev, err
	}
	if after != nil {
		after(&ev, raw)
	}
	return ev, nil
}

// EventDecoder returns a Decoder that unmarshals into T and wraps failures with c.
func EventDecoder[T Event](c *Codec) Decoder {
	return func(raw []byte) (Event, error) {
		return DecodeEvent[T](c, raw, nil)
	}
}

// DecodeEvent unmarshals raw into T, optionally runs after, and wraps failures with c.
func DecodeEvent[T Event](c *Codec, raw []byte, after func(*T, []byte)) (Event, error) {
	ev, err := DecodeAsAndThen(raw, after)
	if err != nil {
		return nil, c.WrapDecodeError(ev, err)
	}
	return ev, nil
}

// PeekHookEventName extracts hook_event_name from a JSON object payload.
func PeekHookEventName(raw []byte) (string, error) {
	var peek struct {
		HookEventName string `json:"hook_event_name"`
	}
	if err := json.Unmarshal(raw, &peek); err != nil {
		return "", err
	}
	return peek.HookEventName, nil
}

// RequireHookEventName peeks hook_event_name and rejects empty stdin or a missing name.
// empty, decodeErr, and nameRequired are the caller's sentinel errors.
func RequireHookEventName(raw []byte, empty, decodeErr, nameRequired error) (string, error) {
	if len(raw) == 0 {
		return "", empty
	}
	name, err := PeekHookEventName(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", decodeErr, err)
	}
	if name == "" {
		return "", nameRequired
	}
	return name, nil
}
