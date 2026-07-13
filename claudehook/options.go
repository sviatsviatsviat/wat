package claudehook

import "os"

// FailPolicy selects exit behavior when a handler returns an error.
type FailPolicy int

const (
	// FailOpen exits with code 1 on handler errors (default).
	FailOpen FailPolicy = iota
	// FailBlock exits with code 2 on handler errors.
	FailBlock
)

// Option configures Mux Serve/Main and Encode side effects.
type Option func(*runtimeConfig)

type runtimeConfig struct {
	policy     FailPolicy
	getenv     func(string) string
	appendFile func(path string, data []byte) error
}

func applyOptions(cfg *runtimeConfig, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
	}
}

func defaultRuntimeConfig() runtimeConfig {
	return runtimeConfig{
		policy: FailOpen,
		getenv: os.Getenv,
	}
}

// WithFailPolicy sets handler error exit behavior.
func WithFailPolicy(p FailPolicy) Option {
	return func(c *runtimeConfig) {
		c.policy = p
	}
}

// WithGetenv injects environment lookup for CLAUDE_ENV_FILE side effects.
func WithGetenv(getenv func(string) string) Option {
	return func(c *runtimeConfig) {
		c.getenv = getenv
	}
}

// WithAppendFile injects file append for CLAUDE_ENV_FILE side effects in tests.
func WithAppendFile(appendFile func(path string, data []byte) error) Option {
	return func(c *runtimeConfig) {
		c.appendFile = appendFile
	}
}

func (c runtimeConfig) encodeOpts() []Option {
	var opts []Option
	if c.getenv != nil {
		opts = append(opts, WithGetenv(c.getenv))
	}
	if c.appendFile != nil {
		opts = append(opts, WithAppendFile(c.appendFile))
	}
	return opts
}
