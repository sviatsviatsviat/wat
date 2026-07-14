// Package claude adapts Claude Code hook payloads to the unified agnostic model.
package claude

import (
	"errors"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Codec implements model.Codec for Claude Code hooks by adapting claude types.
// Reference: https://code.claude.com/docs/en/hooks
type Codec struct {
	// Getenv and AppendFile are injectable for tests. They back the
	// CLAUDE_ENV_FILE side effect used to express Result.Env.
	Getenv     func(string) string
	AppendFile func(path string, data []byte) error
}

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

// Encode renders a unified Result as Claude Code stdout JSON. Claude ignores
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
	var opts []sdkclaude.Option
	if c.Getenv != nil {
		opts = append(opts, sdkclaude.WithGetenv(c.Getenv))
	}
	if c.AppendFile != nil {
		opts = append(opts, sdkclaude.WithAppendFile(c.AppendFile))
	}
	b, err := sdkclaude.Encode(ev.Name, out, opts...)
	if err != nil {
		return nil, 0, fmt.Errorf("claude: %w", err)
	}
	return b, 0, err
}

// RuntimeConfig carries Claude-specific run.Main wiring for agnostic encode.
type RuntimeConfig struct {
	// AppendFile injects file append for CLAUDE_ENV_FILE side effects in tests.
	AppendFile func(path string, data []byte) error
}

// ApplyRunConfig applies run.Main configuration to codec.
func ApplyRunConfig(codec *Codec, cfg *run.Config) {
	if cfg == nil || codec == nil {
		return
	}
	if cfg.Getenv != nil {
		codec.Getenv = cfg.Getenv
	}
	if v := cfg.DialectConfig("claude"); v != nil {
		if rc, ok := v.(*RuntimeConfig); ok && rc != nil && rc.AppendFile != nil {
			codec.AppendFile = rc.AppendFile
		}
	}
}

func mapDecodeError(err error) error {
	switch {
	case errors.Is(err, sdkclaude.ErrDecodePayload), errors.Is(err, sdkclaude.ErrEmptyPayload):
		return adapter.MapDecodeErrorMessage(err, "claude", "claude")
	default:
		return fmt.Errorf("claude: %w", err)
	}
}
