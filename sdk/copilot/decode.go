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

// payloadPeek is a minimal discriminant for format, event name, and Stop scope.
type payloadPeek struct {
	HookEventName         string `json:"hook_event_name"`
	SessionID             string `json:"sessionId"`
	AgentName             string `json:"agent_name"`
	AgentNameCamel        string `json:"agentName"`
	AgentDisplayName      string `json:"agent_display_name"`
	AgentDisplayNameCamel string `json:"agentDisplayName"`
}

func peekPayload(raw []byte) (payloadPeek, error) {
	var peek payloadPeek
	if err := json.Unmarshal(raw, &peek); err != nil {
		return peek, err
	}
	return peek, nil
}

func (p payloadPeek) format() Format {
	if p.HookEventName != "" {
		return FormatVSCode
	}
	if p.SessionID != "" {
		return FormatCamel
	}
	return FormatUnknown
}

func (p payloadPeek) hasSubagentScope() bool {
	return p.AgentName != "" || p.AgentNameCamel != "" ||
		p.AgentDisplayName != "" || p.AgentDisplayNameCamel != ""
}

// SniffFormat detects the Copilot wire format from a payload.
func SniffFormat(raw []byte) Format {
	peek, err := peekPayload(raw)
	if err != nil {
		return FormatUnknown
	}
	return peek.format()
}

// decode parses a hook stdin payload into a typed Event.
func decode(raw []byte, opts ...Option) (Event, error) {
	cfg := defaultDecodeConfig()
	applyOptions(&cfg, opts...)
	return decodeWithHint(raw, cfg.eventHint.Hint)
}

// decodeWithHint parses raw using an explicit event hint.
// It peeks format and event name once, then unmarshals into the matching type.
func decodeWithHint(raw []byte, eventHint string) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}

	peek, err := peekPayload(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if peek.format() == FormatUnknown {
		return nil, ErrUnrecognizedFormat
	}

	received := peek.HookEventName
	if received == "" {
		received = eventHint
	}
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
func eventNameFromRaw(raw []byte, eventHint string) (string, error) {
	if len(raw) == 0 {
		return "", ErrEmptyPayload
	}
	peek, err := peekPayload(raw)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrDecodePayload, err)
	}
	if peek.format() == FormatUnknown {
		return "", ErrUnrecognizedFormat
	}
	received := peek.HookEventName
	if received == "" {
		received = eventHint
	}
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
