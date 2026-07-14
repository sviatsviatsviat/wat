package copilot

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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
		return nil, fmt.Errorf("copilot: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
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
		received = cfg.eventHint.Hint
	}
	if received == "" {
		return nil, ErrEventNameRequired
	}

	canonical, known := ResolveCanonical(raw, received)
	if !known {
		env.setEnvelopeMeta(received, "", raw)
		return RawEvent{Envelope: env, Raw: hookkit.CloneBytes(raw)}, nil
	}

	if fn, ok := decoders[canonical]; ok {
		return fn(raw, received, canonical)
	}
	env.setEnvelopeMeta(received, canonical, raw)
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
