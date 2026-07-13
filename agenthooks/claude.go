package agenthooks

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/claudehook"
)

// ClaudeCodec implements Codec for Claude Code hooks by adapting claudehook types.
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
	native, err := claudehook.Decode(raw)
	if err != nil {
		return nil, fmt.Errorf("claude: %w", err)
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
	var opts []claudehook.Option
	if c.Getenv != nil {
		opts = append(opts, claudehook.WithGetenv(c.Getenv))
	}
	if c.AppendFile != nil {
		opts = append(opts, claudehook.WithAppendFile(c.AppendFile))
	}
	b, err := claudehook.Encode(ev.Name, out, opts...)
	if err != nil {
		return nil, 0, fmt.Errorf("claude: %w", err)
	}
	return b, 0, err
}

// As re-decodes Event.Raw into a claudehook native event type.
func As[T claudehook.Event](ev *Event) (T, error) {
	var zero T
	if ev == nil || len(ev.Raw) == 0 {
		return zero, fmt.Errorf("claude: As: empty event")
	}
	if ev.Agent != Claude {
		return zero, fmt.Errorf("claude: As: event is %v, not Claude", ev.Agent)
	}
	native, err := claudehook.Decode(ev.Raw)
	if err != nil {
		return zero, err
	}
	typed, ok := native.(T)
	if !ok {
		return zero, fmt.Errorf("claude: As: decoded %T, want %T", native, zero)
	}
	return typed, nil
}
