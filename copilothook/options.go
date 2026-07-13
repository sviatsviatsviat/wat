package copilothook

// Option configures Decode, Serve, and Main.
type Option func(*runtimeConfig)

type runtimeConfig struct {
	eventHint string
}

type decodeConfig = runtimeConfig

func defaultDecodeConfig() decodeConfig {
	return decodeConfig{}
}

func applyOptions(cfg *runtimeConfig, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// WithEvent supplies the native event name for camelCase payloads that omit
// hook_event_name. Required for most camelCase CLI payloads.
func WithEvent(name string) Option {
	return func(c *runtimeConfig) {
		c.eventHint = name
	}
}
