package agnostic

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/claude"
)

// ClaudeCodec implements Codec for Claude Code hooks by adapting claude types.
// Reference: https://code.claude.com/docs/en/hooks
type ClaudeCodec struct {
	// Getenv and AppendFile are injectable for tests. They back the
	// CLAUDE_ENV_FILE side effect used to express Result.Env.
	Getenv     func(string) string
	AppendFile func(path string, data []byte) error
}

// Dialect returns Claude.
func (c *ClaudeCodec) Dialect() Dialect { return Claude }

// Decode parses a Claude Code hook stdin payload into a unified Event.
func (c *ClaudeCodec) Decode(raw []byte, eventHint string) (*Event, error) {
	native, err := claude.Decode(raw)
	if err != nil {
		return nil, mapClaudeDecodeError(err)
	}
	ev := mapClaudeEvent(native, raw)
	if ev.Name == "" && eventHint != "" {
		ev.Name = eventHint
		if k, ok := ClaudeKindForEvent(eventHint); ok {
			ev.Kind = k
		}
	}
	return ev, nil
}

// Encode renders a unified Result as Claude Code stdout JSON. Claude ignores
// exit 2 with JSON, so blocking is expressed via fields and exit code is always 0.
// ev must be non-nil.
func (c *ClaudeCodec) Encode(ev *Event, res Result) ([]byte, int, error) {
	if ev == nil {
		return nil, 0, fmt.Errorf("claude: encode: nil event")
	}
	if res.IsZero() {
		return nil, 0, nil
	}
	out := mapClaudeOutput(ev, res)
	var opts []claude.Option
	if c.Getenv != nil {
		opts = append(opts, claude.WithGetenv(c.Getenv))
	}
	if c.AppendFile != nil {
		opts = append(opts, claude.WithAppendFile(c.AppendFile))
	}
	b, err := claude.Encode(ev.Name, out, opts...)
	if err != nil {
		return nil, 0, fmt.Errorf("claude: %w", err)
	}
	return b, 0, err
}

func mapClaudeDecodeError(err error) error {
	switch {
	case errors.Is(err, claude.ErrDecodePayload), errors.Is(err, claude.ErrEmptyPayload):
		return mapDecodeErrorMessage(err, "claude", "claude")
	default:
		return fmt.Errorf("claude: %w", err)
	}
}
