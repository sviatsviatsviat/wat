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
	envelopeAccessorForValue(&ev).envelopePtr().setEnvelopeMeta(received, canonical, raw)
	return ev, nil
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
	var rawEventRaw json.RawMessage
	if e, ok := ev.(RawEvent); ok {
		rawEventRaw = e.Raw
	}
	var accessor hookkit.EnvelopeAccessor
	if rawEventRaw == nil {
		accessor = envelopeAccessorForEvent(ev).envelopePtr()
	}
	return hookkit.RawBytes(ev, rawEventRaw, accessor, func(v any) ([]byte, error) {
		return json.Marshal(v)
	})
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return *envelopeAccessorForEvent(ev).envelopePtr()
}

func decodeWithHint(raw []byte, hint string) (Event, error) {
	return DecodeWithHint(raw, hint)
}
