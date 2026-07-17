package cursor

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// WarnExit is exit code 2. Cursor treats it as block/deny, equivalent to
// returning permission:"deny" on permission-gating events.
const WarnExit = sdkcursor.PermissionDenyExit

// HandlerErrorExit is exit code 1. The runner should use this when a handler
// returns an error under Cursor's default fail-open policy.
const HandlerErrorExit = sdkcursor.HandlerErrorExit

// Decode parses a Cursor hook stdin payload into a unified Event.
func Decode(raw []byte, eventHint string) (*model.Event, error) {
	native, err := sdkcursor.Decode(raw, sdkcursor.WithEvent(eventHint))
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
	if ev.Name == "" {
		return nil, mapDecodeError(sdkcursor.ErrEventNameRequired)
	}
	return ev, nil
}

func mapDecodeError(err error) error {
	switch {
	case errors.Is(err, sdkcursor.ErrDecodePayload), errors.Is(err, sdkcursor.ErrEmptyPayload):
		return adapter.MapDecodeErrorMessage(err, "cursor", "cursor")
	case errors.Is(err, sdkcursor.ErrEventNameRequired):
		return adapter.RemapDecodeError(err, "cursor: decode: event name required (use eventHint)")
	default:
		return fmt.Errorf("cursor: %w", err)
	}
}
