package claude

import (
	"os"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// FailPolicy selects exit behavior when a handler returns an error.
type FailPolicy int

const (
	// FailOpen exits with code 1 on handler errors (default).
	FailOpen FailPolicy = iota
	// FailBlock exits with code 2 on handler errors.
	FailBlock
)

// Option configures Encode side effects.
type Option func(*runtimeConfig)

type runtimeConfig struct {
	policy     FailPolicy
	getenv     func(string) string
	appendFile func(path string, data []byte) error
}

func defaultRuntimeConfig() runtimeConfig {
	return runtimeConfig{
		policy: FailOpen,
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
	if v := cfg.DialectConfig("claude"); v != nil {
		if rc, ok := v.(*runtimeConfig); ok && rc != nil {
			return rc
		}
	}
	rc := defaultRuntimeConfig()
	cfg.SetDialectConfig("claude", &rc)
	return &rc
}

// WithFailPolicy sets handler error exit behavior for claude handlers.
func WithFailPolicy(p FailPolicy) run.Option {
	return func(c *run.Config) {
		rc := claudeRunConfig(c)
		rc.policy = p
	}
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

// WithRunGetenv injects environment lookup for claude handlers at serve time.
func WithRunGetenv(getenv func(string) string) run.Option {
	return func(c *run.Config) {
		rc := claudeRunConfig(c)
		rc.getenv = getenv
	}
}

// WithRunAppendFile injects file append for claude handlers at serve time.
func WithRunAppendFile(appendFile func(path string, data []byte) error) run.Option {
	return func(c *run.Config) {
		rc := claudeRunConfig(c)
		rc.appendFile = appendFile
	}
}
