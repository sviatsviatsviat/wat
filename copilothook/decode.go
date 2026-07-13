package copilothook

import (
	"encoding/json"
	"fmt"
	"io"
)

type decodeFn func([]byte, string, string) (Event, error)

var decoders = map[string]decodeFn{
	EventSessionStart:        decodeAs[SessionStart],
	EventSessionEnd:          decodeAs[SessionEnd],
	EventUserPromptSubmitted: decodeAs[UserPromptSubmitted],
	EventPreToolUse:          decodeAs[PreToolUse],
	EventPostToolUse:         decodeAs[PostToolUse],
	EventPostToolUseFailure:  decodeAs[PostToolUseFailure],
	EventPermissionRequest:   decodeAs[PermissionRequest],
	EventSubagentStart:       decodeAs[SubagentStart],
	EventSubagentStop:        decodeAs[SubagentStop],
	EventAgentStop:           decodeAs[AgentStop],
	EventPreCompact:          decodeAs[PreCompact],
	EventNotification:        decodeAs[Notification],
	EventErrorOccurred:       decodeAs[ErrorOccurred],
}

func decodeAs[T Event](raw []byte, received, canonical string) (Event, error) {
	var ev T
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, fmt.Errorf("copilothook: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
	}
	envelopeAccessorForValue(&ev).envelopePtr().setEnvelopeMeta(received, canonical, raw)
	return ev, nil
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

// Decode parses a GitHub Copilot hook stdin payload into a typed Event.
// camelCase payloads require WithEvent unless hook_event_name is present.
func Decode(raw []byte, opts ...Option) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}
	cfg := defaultDecodeConfig()
	applyOptions(&cfg, opts...)

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
		received = cfg.eventHint
	}
	if received == "" {
		return nil, ErrEventNameRequired
	}

	canonical, known := ResolveCanonical(raw, received)
	if !known {
		env.setEnvelopeMeta(received, "", raw)
		return RawEvent{Envelope: env, Raw: cloneRaw(raw)}, nil
	}

	if fn, ok := decoders[canonical]; ok {
		return fn(raw, received, canonical)
	}
	env.setEnvelopeMeta(received, canonical, raw)
	return RawEvent{Envelope: env, Raw: cloneRaw(raw)}, nil
}

// ParseEvent reads and decodes a hook payload from r.
func ParseEvent(r io.Reader, opts ...Option) (Event, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("copilothook: read payload: %w", err)
	}
	return Decode(raw, opts...)
}

func cloneRaw(raw []byte) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), raw...)
}

// RawBytes returns the untouched JSON for an event when available.
func RawBytes(ev Event) json.RawMessage {
	switch e := ev.(type) {
	case RawEvent:
		return cloneRaw(e.Raw)
	default:
		if raw := decodedRawFrom(ev); len(raw) > 0 {
			return cloneRaw(raw)
		}
		b, err := json.Marshal(ev)
		if err != nil {
			return nil
		}
		return b
	}
}

func decodedRawFrom(ev Event) json.RawMessage {
	return envelopeAccessorForEvent(ev).envelopePtr().decodedRawBytes()
}

// CloneRaw copies raw JSON for independent mutation.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return cloneRaw(raw)
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return *envelopeAccessorForEvent(ev).envelopePtr()
}
