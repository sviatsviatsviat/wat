package claude

import (
	"os"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Option configures Encode side effects.
type Option func(*runtimeConfig)

type runtimeConfig struct {
	getenv     func(string) string
	appendFile func(path string, data []byte) error
}

func defaultRuntimeConfig() runtimeConfig {
	return runtimeConfig{
		getenv: os.Getenv,
	}
}

func applyOptions(cfg *runtimeConfig, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
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

func claudeRunConfig(cfg *run.Config) *runtimeConfig {
	if v := cfg.DialectConfig(Dialect); v != nil {
		if rc, ok := v.(*runtimeConfig); ok && rc != nil {
			return rc
		}
	}
	rc := defaultRuntimeConfig()
	cfg.SetDialectConfig(Dialect, &rc)
	return &rc
}

// WithGetenv injects environment lookup for CLAUDE_ENV_FILE Encode side effects.
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
