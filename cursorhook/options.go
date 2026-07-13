package cursorhook

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

// WithEvent supplies the native event name when hook_event_name is absent.
func WithEvent(name string) Option {
	return func(c *runtimeConfig) {
		c.eventHint = name
	}
}
