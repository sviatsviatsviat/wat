package agenthooks

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/cursorhook"
)

// CursorWarnExit is exit code 2. Cursor treats it as block/deny, equivalent to
// returning permission:"deny" on permission-gating events.
const CursorWarnExit = cursorhook.PermissionDenyExit

// CursorHandlerErrorExit is exit code 1. The runner should use this when a
// handler returns an error under Cursor's default fail-open policy.
const CursorHandlerErrorExit = cursorhook.HandlerErrorExit

// CursorCodec implements Codec for Cursor hooks by adapting cursorhook types.
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
	native, err := cursorhook.Decode(raw, cursorhook.WithEvent(eventHint))
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
	b, code, err := cursorhook.Encode(cursorEncodeEventName(ev), out)
	if err != nil {
		return nil, 0, fmt.Errorf("cursor: %w", err)
	}
	return b, code, err
}

// AsCursor re-decodes Event.Raw into a cursorhook native event type.
func AsCursor[T cursorhook.Event](ev *Event) (T, error) {
	var zero T
	if ev == nil || len(ev.Raw) == 0 {
		return zero, fmt.Errorf("cursor: AsCursor: empty event")
	}
	if ev.Agent != Cursor {
		return zero, fmt.Errorf("cursor: AsCursor: event is %v, not Cursor", ev.Agent)
	}
	native, err := cursorhook.Decode(ev.Raw, cursorhook.WithEvent(ev.Name))
	if err != nil {
		return zero, err
	}
	typed, ok := native.(T)
	if !ok {
		return zero, fmt.Errorf("cursor: AsCursor: decoded %T, want %T", native, zero)
	}
	return typed, nil
}

func mapCursorDecodeError(err error) error {
	switch {
	case errors.Is(err, cursorhook.ErrDecodePayload), errors.Is(err, cursorhook.ErrEmptyPayload):
		return mapDecodeErrorMessage(err, "cursor", "cursorhook")
	case errors.Is(err, cursorhook.ErrEventNameRequired):
		return remapDecodeError(err, "cursor: decode: event name required (use eventHint)")
	default:
		return fmt.Errorf("cursor: %w", err)
	}
}
