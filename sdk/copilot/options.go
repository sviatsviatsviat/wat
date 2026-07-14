package copilot

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Option configures Decode.
type Option func(*decodeConfig)

type decodeConfig struct {
	eventHint hookkit.EventHint
}

func defaultDecodeConfig() decodeConfig {
	return decodeConfig{eventHint: hookkit.DefaultEventHint()}
}

func applyOptions(cfg *decodeConfig, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// WithEvent supplies the native event name for camelCase payloads that omit
// hook_event_name. Required for most camelCase CLI payloads when calling Decode.
func WithEvent(name string) Option {
	return func(c *decodeConfig) {
		hookkit.ApplyHintOptions(&c.eventHint, hookkit.WithEventHint(name))
	}
}
