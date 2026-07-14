// Package claude adapts Claude Code hook payloads to the unified agnostic model.
package claude

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// Codec implements model.Codec for Claude Code hooks by adapting claude types.
// Reference: https://code.claude.com/docs/en/hooks
type Codec struct{}

// Dialect returns Claude.
func (c *Codec) Dialect() model.Dialect { return model.Claude }

// Decode parses a Claude Code hook stdin payload into a unified Event.
func (c *Codec) Decode(raw []byte, eventHint string) (*model.Event, error) {
	native, err := sdkclaude.Decode(raw)
	if err != nil {
		return nil, mapDecodeError(err)
	}
	ev := mapEvent(native, raw)
	if ev.Name == "" && eventHint != "" {
		ev.Name = eventHint
		if k, ok := KindForEvent(eventHint); ok {
			ev.Kind = k
		}
	}
	return ev, nil
}

// Encode renders a portable Result as Claude Code stdout JSON. Claude ignores
// exit 2 with JSON, so blocking is expressed via fields and exit code is always 0.
// ev must be non-nil.
func (c *Codec) Encode(ev *model.Event, res model.Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("claude: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := mapOutput(ev, res)
	if out == nil {
		return nil, 0, fmt.Errorf("claude: encode: %s has no portable encode surface", ev.Kind)
	}
	b, err := sdkclaude.Encode(ev.Name, out)
	if err != nil {
		return nil, 0, fmt.Errorf("claude: %w", err)
	}
	return b, 0, err
}

func mapDecodeError(err error) error {
	switch {
	case errors.Is(err, sdkclaude.ErrDecodePayload), errors.Is(err, sdkclaude.ErrEmptyPayload):
		return adapter.MapDecodeErrorMessage(err, "claude", "claude")
	default:
		return fmt.Errorf("claude: %w", err)
	}
}
