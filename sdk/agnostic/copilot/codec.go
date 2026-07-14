// Package copilot adapts GitHub Copilot hook payloads to the unified agnostic model.
package copilot

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// PreToolErrorExit is the exit code when a preToolUse handler returns an error.
// Copilot command hooks fail-closed on non-zero exits other than 2.
const PreToolErrorExit = sdkcopilot.PreToolErrorExit

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// permissionRequest it means deny, and for postToolUseFailure it carries
// additionalContext in stdout.
const WarnExit = sdkcopilot.WarnExit

// Codec implements model.Codec for GitHub Copilot hooks by adapting copilot types.
//
// Handler errors on preToolUse should exit PreToolErrorExit (fail-closed).
// Encode returns WarnExit only for documented output paths on permissionRequest
// deny and postToolUseFailure context.
//
// Reference: https://docs.github.com/en/copilot/reference/hooks-reference
type Codec struct{}

// Dialect returns Copilot.
func (c *Codec) Dialect() model.Dialect { return model.Copilot }

// Decode parses a GitHub Copilot hook stdin payload into a unified Event.
func (c *Codec) Decode(raw []byte, eventHint string) (*model.Event, error) {
	native, err := sdkcopilot.Decode(raw, sdkcopilot.WithEvent(eventHint))
	if err != nil {
		return nil, mapDecodeError(err)
	}
	ev := MapEvent(native, raw)
	if ev.Name == "" && eventHint != "" {
		ev.Name = eventHint
		if k, ok := KindForEvent(eventHint); ok {
			ev.Kind = k
		}
	}
	return ev, nil
}

// Encode renders a unified Result as Copilot stdout JSON and exit code.
// ev must be non-nil.
func (c *Codec) Encode(ev *model.Event, res model.Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("copilot: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := MapOutput(ev, res)
	if out == nil {
		return nil, 0, fmt.Errorf("copilot: encode: %s has no portable encode surface", ev.Kind)
	}
	b, code, err := sdkcopilot.Encode(encodeEventName(ev), out)
	if err != nil {
		return nil, 0, fmt.Errorf("copilot: %w", err)
	}
	return b, code, err
}

func mapDecodeError(err error) error {
	switch {
	case errors.Is(err, sdkcopilot.ErrUnrecognizedFormat):
		return adapter.RemapDecodeError(err, "copilot: decode payload: unrecognized format")
	case errors.Is(err, sdkcopilot.ErrEventNameRequired):
		return adapter.RemapDecodeError(err, "copilot: decode: event name required (camelCase payloads need eventHint)")
	case errors.Is(err, sdkcopilot.ErrDecodePayload), errors.Is(err, sdkcopilot.ErrEmptyPayload):
		return adapter.MapDecodeErrorMessage(err, "copilot", "copilot")
	default:
		return fmt.Errorf("copilot: %w", err)
	}
}

func encodeEventName(ev *model.Event) string {
	if ev.Name != "" {
		return ev.Name
	}
	return EventForKind[ev.Kind]
}
