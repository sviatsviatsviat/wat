package agnostic

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/cursor"
)

// CursorWarnExit is exit code 2. Cursor treats it as block/deny, equivalent to
// returning permission:"deny" on permission-gating events.
const CursorWarnExit = cursor.PermissionDenyExit

// CursorHandlerErrorExit is exit code 1. The runner should use this when a
// handler returns an error under Cursor's default fail-open policy.
const CursorHandlerErrorExit = cursor.HandlerErrorExit

// CursorCodec implements Codec for Cursor hooks by adapting cursor types.
//
// Cursor's dedicated surfaces are folded into the unified tool kinds:
// beforeShellExecution/beforeMCPExecution/beforeReadFile → KindPreTool,
// afterShellExecution/afterMCPExecution/afterFileEdit → KindPostTool.
// Event.Name preserves the native surface; Event.Raw holds the full payload.
//
// Reference: https://cursor.com/docs/hooks
type CursorCodec struct{}

// Dialect returns Cursor.
func (c *CursorCodec) Dialect() Dialect { return Cursor }

// Decode parses a Cursor hook stdin payload into a unified Event.
func (c *CursorCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	native, err := cursor.Decode(raw, cursor.WithEvent(eventHint))
	if err != nil {
		return nil, mapCursorDecodeError(err)
	}
	ev := mapCursorEvent(native, raw)
	if ev.Name == "" && eventHint != "" {
		ev.Name = eventHint
		if k, ok := CursorKindForEvent(eventHint); ok {
			ev.Kind = k
		}
	}
	if ev.Name == "" {
		return nil, mapCursorDecodeError(cursor.ErrEventNameRequired)
	}
	return ev, nil
}

// Encode renders a unified Result as Cursor stdout JSON and exit code.
// ev must be non-nil.
func (c *CursorCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("cursor: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := mapCursorOutput(ev, res)
	if out == nil {
		return nil, 0, nil
	}
	b, code, err := cursor.Encode(cursorEncodeEventName(ev), out)
	if err != nil {
		return nil, 0, fmt.Errorf("cursor: %w", err)
	}
	return b, code, err
}

func mapCursorDecodeError(err error) error {
	switch {
	case errors.Is(err, cursor.ErrDecodePayload), errors.Is(err, cursor.ErrEmptyPayload):
		return mapDecodeErrorMessage(err, "cursor", "cursor")
	case errors.Is(err, cursor.ErrEventNameRequired):
		return remapDecodeError(err, "cursor: decode: event name required (use eventHint)")
	default:
		return fmt.Errorf("cursor: %w", err)
	}
}
