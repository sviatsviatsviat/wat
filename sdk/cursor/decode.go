package cursor

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sviatsviatsviat/wat/sdk/internal/hookkit"
)

type decodeFn func([]byte, string, string) (Event, error)

var decoders = map[string]decodeFn{
	EventSessionStart:         decodeAs[SessionStart],
	EventSessionEnd:           decodeAs[SessionEnd],
	EventBeforeSubmitPrompt:   decodeAs[BeforeSubmitPrompt],
	EventPreToolUse:           decodeAs[PreToolUse],
	EventPostToolUse:          decodeAs[PostToolUse],
	EventPostToolUseFailure:   decodeAs[PostToolUseFailure],
	EventBeforeShellExecution: decodeAs[BeforeShellExecution],
	EventAfterShellExecution:  decodeAs[AfterShellExecution],
	EventBeforeMCPExecution:   decodeAs[BeforeMCPExecution],
	EventAfterMCPExecution:    decodeAs[AfterMCPExecution],
	EventBeforeReadFile:       decodeAs[BeforeReadFile],
	EventAfterFileEdit:        decodeAs[AfterFileEdit],
	EventSubagentStart:        decodeAs[SubagentStart],
	EventSubagentStop:         decodeAs[SubagentStop],
	EventStop:                 decodeAs[Stop],
	EventPreCompact:           decodeAs[PreCompact],
	EventAfterAgentResponse:   decodeAs[AfterAgentResponse],
	EventAfterAgentThought:    decodeAs[AfterAgentThought],
	EventBeforeTabFileRead:    decodeAs[BeforeTabFileRead],
	EventAfterTabFileEdit:     decodeAs[AfterTabFileEdit],
	EventWorkspaceOpen:        decodeAs[WorkspaceOpen],
}

func decodeAs[T Event](raw []byte, received, canonical string) (Event, error) {
	var ev T
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, fmt.Errorf("cursor: decode %T: %w", ev, fmt.Errorf("%w: %w", ErrDecodePayload, err))
	}
	envelopeAccessorForValue(&ev).envelopePtr().setEnvelopeMeta(received, canonical, raw)
	return ev, nil
}

// Decode parses a Cursor hook stdin payload into a typed Event.
func Decode(raw []byte, opts ...Option) (Event, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPayload
	}
	cfg := defaultDecodeConfig()
	applyOptions(&cfg, opts...)

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

	canonical := received
	if fn, ok := decoders[canonical]; ok {
		return fn(raw, received, canonical)
	}
	env.setEnvelopeMeta(received, "", raw)
	return RawEvent{Envelope: env, Raw: hookkit.CloneBytes(raw)}, nil
}

// ParseEvent reads and decodes a hook payload from r.
func ParseEvent(r io.Reader, opts ...Option) (Event, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cursor: read payload: %w", err)
	}
	return Decode(raw, opts...)
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

// CloneRaw copies raw JSON for independent mutation.
func CloneRaw(raw json.RawMessage) json.RawMessage {
	return hookkit.CloneRaw(raw)
}

// EnvelopeOf returns the shared envelope from a decoded event.
func EnvelopeOf(ev Event) Envelope {
	return *envelopeAccessorForEvent(ev).envelopePtr()
}
