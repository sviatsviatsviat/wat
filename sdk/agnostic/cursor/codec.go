// Package cursor adapts Cursor hook payloads to the unified agnostic model.
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

// Codec implements model.Codec for Cursor hooks by adapting cursor types.
//
// Cursor's dedicated surfaces are folded into the unified tool kinds:
// beforeShellExecution/beforeMCPExecution/beforeReadFile → KindPreTool,
// afterShellExecution/afterMCPExecution/afterFileEdit → KindPostTool.
// Event.Name preserves the native surface; Event.Raw holds the full payload.
//
// Reference: https://cursor.com/docs/hooks
type Codec struct{}

// Dialect returns Cursor.
func (c *Codec) Dialect() model.Dialect { return model.Cursor }

// Decode parses a Cursor hook stdin payload into a unified Event.
func (c *Codec) Decode(raw []byte, eventHint string) (*model.Event, error) {
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

// Encode renders a unified Result as Cursor stdout JSON and exit code.
// ev must be non-nil.
func (c *Codec) Encode(ev *model.Event, res model.Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("cursor: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := MapOutput(ev, res)
	if out == nil {
		return nil, 0, fmt.Errorf("cursor: encode: %s has no portable encode surface", ev.Kind)
	}
	b, code, err := sdkcursor.Encode(encodeEventName(ev), out)
	if err != nil {
		return nil, 0, fmt.Errorf("cursor: %w", err)
	}
	return b, code, err
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

func encodeEventName(ev *model.Event) string {
	if ev.Name != "" {
		return ev.Name
	}
	return EventForKind[ev.Kind]
}
