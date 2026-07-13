package agnostic

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
)

// CopilotPreToolErrorExit is the exit code when a preToolUse handler returns
// an error. Copilot command hooks fail-closed on non-zero exits other than 2.
const CopilotPreToolErrorExit = copilot.PreToolErrorExit

// CopilotHandlerErrorExit is exit code 1 for handler errors under fail-open policy.
const CopilotHandlerErrorExit = copilot.HandlerErrorExit

// CopilotWarnExit is exit code 2. Copilot treats it as a warning by default;
// for permissionRequest it means deny, and for postToolUseFailure it carries
// additionalContext in stdout.
const CopilotWarnExit = copilot.WarnExit

// CopilotCodec implements Codec for GitHub Copilot hooks by adapting copilot types.
//
// Handler errors on preToolUse should exit CopilotPreToolErrorExit (fail-closed).
// Encode returns CopilotWarnExit only for documented output paths on
// permissionRequest deny and postToolUseFailure context.
//
// Reference: https://docs.github.com/en/copilot/reference/hooks-reference
type CopilotCodec struct{}

// Dialect returns Copilot.
func (c *CopilotCodec) Dialect() Dialect { return Copilot }

// Decode parses a GitHub Copilot hook stdin payload into a unified Event.
func (c *CopilotCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	native, err := copilot.Decode(raw, copilot.WithEvent(eventHint))
	if err != nil {
		return nil, mapCopilotDecodeError(err)
	}
	ev := mapCopilotEvent(native, raw)
	if ev.Name == "" && eventHint != "" {
		ev.Name = eventHint
		if k, ok := CopilotKindForEvent(eventHint); ok {
			ev.Kind = k
		}
	}
	return ev, nil
}

// Encode renders a unified Result as Copilot stdout JSON and exit code.
// ev must be non-nil.
func (c *CopilotCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("copilot: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := mapCopilotOutput(ev, res)
	if out == nil {
		return nil, 0, nil
	}
	b, code, err := copilot.Encode(copilotEncodeEventName(ev), out)
	if err != nil {
		return nil, 0, fmt.Errorf("copilot: %w", err)
	}
	return b, code, err
}

func mapCopilotDecodeError(err error) error {
	switch {
	case errors.Is(err, copilot.ErrUnrecognizedFormat):
		return remapDecodeError(err, "copilot: decode payload: unrecognized format")
	case errors.Is(err, copilot.ErrEventNameRequired):
		return remapDecodeError(err, "copilot: decode: event name required (camelCase payloads need eventHint)")
	case errors.Is(err, copilot.ErrDecodePayload), errors.Is(err, copilot.ErrEmptyPayload):
		return mapDecodeErrorMessage(err, "copilot", "copilot")
	default:
		return fmt.Errorf("copilot: %w", err)
	}
}

// AsCopilot re-decodes Event.Raw into a copilot native event type.
func AsCopilot[T copilot.Event](ev *Event) (T, error) {
	var zero T
	if ev == nil || len(ev.Raw) == 0 {
		return zero, fmt.Errorf("copilot: AsCopilot: empty event")
	}
	if ev.Agent != Copilot {
		return zero, fmt.Errorf("copilot: AsCopilot: event is %v, not Copilot", ev.Agent)
	}
	native, err := copilot.Decode(ev.Raw, copilot.WithEvent(ev.Name))
	if err != nil {
		return zero, err
	}
	typed, ok := native.(T)
	if !ok {
		return zero, fmt.Errorf("copilot: AsCopilot: decoded %T, want %T", native, zero)
	}
	return typed, nil
}

func copilotEncodeEventName(ev *Event) string {
	if ev.Name != "" {
		return ev.Name
	}
	return CopilotEventForKind[ev.Kind]
}
