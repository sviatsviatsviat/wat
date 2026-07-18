package claude

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

var decoders = hookkit.NewDecoderRegistry()

type decodeFn func([]byte) (Event, error)

func registerDecoder(name string, fn decodeFn) {
	decoders.Register(name, func(raw []byte, _, _ string) (any, error) {
		return fn(raw)
	})
}

func decodeAs[T Event](raw []byte) (Event, error) {
	return decodeAsAndThen[T](raw, nil)
}

func decodeAsAndThen[T Event](raw []byte, after func(*T, []byte)) (Event, error) {
	ev, err := hookkit.DecodeAsAndThen(raw, after)
	if err != nil {
		return nil, fmt.Errorf("claude: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
	}
	return ev, nil
}

// decode parses a Claude Code hook stdin payload into a typed Event.
// It peeks hook_event_name once, then unmarshals the payload into the matching type.
func decode(raw []byte) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}
	name, err := hookkit.PeekHookEventName(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if name == "" {
		return decodeAs[RawEvent](raw)
	}
	if fn, ok := decoders.Lookup(name); ok {
		ev, err := fn(raw, name, name)
		if err != nil {
			return nil, err
		}
		return ev.(Event), nil
	}
	return decodeAs[RawEvent](raw)
}

// eventNameFromRaw peeks the hook event name without a full typed decode.
func eventNameFromRaw(raw []byte, eventHint string) (string, error) {
	if len(raw) == 0 {
		return "", ErrEmptyPayload
	}
	name, err := hookkit.PeekHookEventName(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if name == "" {
		name = eventHint
	}
	if name == "" {
		return "", fmt.Errorf("claude: empty event name")
	}
	return name, nil
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}
