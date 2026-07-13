package cursor

import "github.com/sviatsviatsviat/wat/sdk/internal/hookkit"

// Option configures Decode, Serve, and Main.
type Option func(*runtimeConfig)

type runtimeConfig struct {
	eventHint hookkit.EventHint
}

type decodeConfig = runtimeConfig

func defaultDecodeConfig() decodeConfig {
	return decodeConfig{eventHint: hookkit.DefaultEventHint()}
}

func applyOptions(cfg *runtimeConfig, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// WithEvent supplies the native event name when hook_event_name is absent.
func WithEvent(name string) Option {
	return func(c *runtimeConfig) {
		hookkit.ApplyHintOptions(&c.eventHint, hookkit.WithEventHint(name))
	}
}
