package hookkit

import "encoding/json"

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
