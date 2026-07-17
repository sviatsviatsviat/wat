package claude

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// Decode parses a Claude Code hook stdin payload into a unified Event.
func Decode(raw []byte, eventHint string) (*model.Event, error) {
	native, err := sdkclaude.Decode(raw)
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
	case errors.Is(err, sdkclaude.ErrDecodePayload), errors.Is(err, sdkclaude.ErrEmptyPayload):
		return adapter.MapDecodeErrorMessage(err, "claude", "claude")
	default:
		return fmt.Errorf("claude: %w", err)
	}
}
