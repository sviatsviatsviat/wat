package cursor

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

var decoders = hookkit.NewDecoderRegistry()

type decodeFn func([]byte, string, string) (Event, error)

func registerDecoder(name string, fn decodeFn) {
	decoders.Register(name, func(raw []byte, received, canonical string) (any, error) {
		return fn(raw, received, canonical)
	})
}

func decodeAs[T Event](raw []byte, received, canonical string) (Event, error) {
	return decodeAsAndThen[T](raw, received, canonical, nil)
}

func decodeAsAndThen[T Event](raw []byte, received, canonical string, after func(*T, []byte)) (Event, error) {
	ev, err := hookkit.DecodeAsAndThen(raw, after)
	if err != nil {
		return nil, fmt.Errorf("cursor: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
	}
	attachEnvelopeMeta(&ev, received, canonical)
	return ev, nil
}

func attachEnvelopeMeta(ev any, received, canonical string) {
	s, ok := ev.(interface {
		setEnvelopeMeta(received, canonical string)
	})
	if !ok {
		panic(fmt.Sprintf("cursor: %T cannot attach envelope meta", ev))
	}
	s.setEnvelopeMeta(received, canonical)
}

// decode parses a Cursor hook stdin payload into a typed Event.
// It peeks hook_event_name once, then unmarshals into the matching type.
func decode(raw []byte) (any, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}

	received, err := hookkit.PeekHookEventName(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if received == "" {
		return nil, ErrEventNameRequired
	}

	canonical := received
	if fn, ok := decoders.Lookup(canonical); ok {
		ev, err := fn(raw, received, canonical)
		if err != nil {
			return nil, err
		}
		return ev.(Event), nil
	}
	return decodeAs[RawEvent](raw, received, "")
}

// eventNameFromRaw peeks the hook event name without a full typed decode.
func eventNameFromRaw(raw []byte) (string, error) {
	if len(raw) == 0 {
		return "", ErrEmptyPayload
	}
	received, err := hookkit.PeekHookEventName(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if received == "" {
		return "", ErrEventNameRequired
	}
	return received, nil
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}
