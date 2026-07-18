package claude

import (
	"encoding/json"
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
	attachDecodedRaw(&ev, raw)
	return ev, nil
}

func attachDecodedRaw(ev any, raw []byte) {
	s, ok := ev.(interface{ setDecodedRaw(json.RawMessage) })
	if !ok {
		panic(fmt.Sprintf("claude: %T cannot attach decoded raw", ev))
	}
	s.setDecodedRaw(raw)
}

// Decode parses a Claude Code hook stdin payload into a typed Event.
func Decode(raw []byte) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	name := env.HookEventName
	if name == "" {
		env.setDecodedRaw(raw)
		return RawEvent{Envelope: env, Raw: hookkit.CloneBytes(raw)}, nil
	}
	if fn, ok := decoders.Lookup(name); ok {
		ev, err := fn(raw, name, name)
		if err != nil {
			return nil, err
		}
		return ev.(Event), nil
	}
	env.setDecodedRaw(raw)
	return RawEvent{Envelope: env, Raw: hookkit.CloneBytes(raw)}, nil
}

// RawBytes returns the untouched JSON for an event when available.
func RawBytes(ev Event) json.RawMessage {
	return hookkit.RawBytes(ev, ev.envelope())
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}
