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

// payloadPeek is a minimal discriminant for event name and Stop scope.
type payloadPeek struct {
	HookEventName    string `json:"hook_event_name"`
	AgentName        string `json:"agent_name"`
	AgentDisplayName string `json:"agent_display_name"`
}

func peekPayload(raw []byte) (payloadPeek, error) {
	var peek payloadPeek
	if err := json.Unmarshal(raw, &peek); err != nil {
		return peek, err
	}
	return peek, nil
}

func (p payloadPeek) hasSubagentScope() bool {
	return p.AgentName != "" || p.AgentDisplayName != ""
}

// decode parses a hook stdin payload into a typed Event.
// It requires hook_event_name, then unmarshals into the matching type.
func decode(raw []byte) (any, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}

	peek, err := peekPayload(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	received := peek.HookEventName
	if received == "" {
		return nil, ErrEventNameRequired
	}

	canonical, known := resolveCanonical(received, peek.hasSubagentScope())
	if !known {
		return decodeAs[RawEvent](raw, received, "")
	}

	if fn, ok := decoders.Lookup(canonical); ok {
		ev, err := fn(raw, received, canonical)
		if err != nil {
			return nil, err
		}
		return ev.(Event), nil
	}
	return decodeAs[RawEvent](raw, received, canonical)
}

// eventNameFromRaw peeks the canonical hook event name without a full typed decode.
func eventNameFromRaw(raw []byte) (string, error) {
	if len(raw) == 0 {
		return "", ErrEmptyPayload
	}
	peek, err := peekPayload(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	received := peek.HookEventName
	if received == "" {
		return "", ErrEventNameRequired
	}
	canonical, known := resolveCanonical(received, peek.hasSubagentScope())
	if !known {
		return received, nil
	}
	return canonical, nil
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return ev.envelope()
}
