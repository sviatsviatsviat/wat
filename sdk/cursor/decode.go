package cursor

import (
	"encoding/json"
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
	attachEnvelopeMeta(&ev, received, canonical, raw)
	return ev, nil
}

func attachEnvelopeMeta(ev any, received, canonical string, raw []byte) {
	s, ok := ev.(interface {
		setEnvelopeMeta(received, canonical string, raw json.RawMessage)
	})
	if !ok {
		panic(fmt.Sprintf("cursor: %T cannot attach envelope meta", ev))
	}
	s.setEnvelopeMeta(received, canonical, raw)
}

// Decode parses a Cursor hook stdin payload into a typed Event.
func Decode(raw []byte, opts ...Option) (Event, error) {
	cfg := defaultDecodeConfig()
	applyOptions(&cfg, opts...)
	return DecodeWithHint(raw, cfg.eventHint.Hint)
}

// DecodeWithHint parses raw using an explicit event hint.
func DecodeWithHint(raw []byte, eventHint string) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}

	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}

	received := env.HookEventName
	if received == "" {
		received = eventHint
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
	env.setEnvelopeMeta(received, "", raw)
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

func decodeWithHint(raw []byte, hint string) (Event, error) {
	return DecodeWithHint(raw, hint)
}
