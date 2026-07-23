package run

import (
	"os"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Option configures Main.
type Option func(*Config)

// WithDialect forces a dialect name ("claude", "copilot", "cursor").
func WithDialect(name string) Option {
	return func(c *Config) {
		c.Dialect = name
	}
}

// WithGetenv injects environment lookup for dialect detection and encode side effects.
func WithGetenv(getenv func(string) string) Option {
	return func(c *Config) {
		c.Getenv = getenv
	}
}

func applyOptions(opts ...Option) *Config {
	cfg := hookkit.DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.Getenv == nil {
		cfg.Getenv = os.Getenv
	}
	applyEnv(cfg)
	return cfg
}

func applyEnv(cfg *Config) {
	if cfg.Dialect == "" {
		if v := cfg.Getenv("WAT_AGENT"); v != "" {
			cfg.Dialect = parseDialectEnv(v)
		}
	}
}

func parseDialectEnv(v string) string {
	switch v {
	case "claude", "claude-code", "claudecode":
		return "claude"
	case "copilot", "github-copilot", "gh":
		return "copilot"
	case "cursor":
		return "cursor"
	default:
		return v
	}
}
