package copilot

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
		return nil, fmt.Errorf("copilot: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
	}
	attachEnvelopeMeta(&ev, received, canonical)
	return ev, nil
}

func attachEnvelopeMeta(ev any, received, canonical string) {
	s, ok := ev.(interface {
		setEnvelopeMeta(received, canonical string)
	})
	if !ok {
		panic(fmt.Sprintf("copilot: %T cannot attach envelope meta", ev))
	}
	s.setEnvelopeMeta(received, canonical)
}

func newRawEvent(env Envelope, received, canonical string, raw []byte) RawEvent {
	ev := RawEvent{Envelope: env}
	ev.setEnvelopeMeta(received, canonical)
	ev.SetRaw(raw)
	return ev
}

// SniffFormat detects the Copilot wire format from a payload.
func SniffFormat(raw []byte) Format {
	var peek struct {
		HookEventName string `json:"hook_event_name"`
		SessionID     string `json:"sessionId"`
	}
	if json.Unmarshal(raw, &peek) != nil {
		return FormatUnknown
	}
	if peek.HookEventName != "" {
		return FormatVSCode
	}
	if peek.SessionID != "" {
		return FormatCamel
	}
	return FormatUnknown
}

// Decode parses a hook stdin payload into a typed Event.
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

	format := SniffFormat(raw)
	if format == FormatUnknown {
		return nil, ErrUnrecognizedFormat
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

	canonical, known := ResolveCanonical(raw, received)
	if !known {
		return newRawEvent(env, received, "", raw), nil
	}

	if fn, ok := decoders.Lookup(canonical); ok {
		ev, err := fn(raw, received, canonical)
		if err != nil {
			return nil, err
		}
		return ev.(Event), nil
	}
	return newRawEvent(env, received, canonical, raw), nil
}

// RawBytes returns the untouched JSON for an event when available.
func RawBytes(ev Event) json.RawMessage {
	return ev.Raw()
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}

func decodeWithHint(raw []byte, hint string) (Event, error) {
	return DecodeWithHint(raw, hint)
}
