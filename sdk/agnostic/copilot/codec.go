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

// Decode parses a GitHub Copilot hook stdin payload into a unified Event.
func Decode(raw []byte, eventHint string) (*model.Event, error) {
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
