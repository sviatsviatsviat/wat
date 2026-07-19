package run

import (
	"context"
	"os"
)

// Option configures Main.
type Option func(*Config)

// Config holds resolved serve-time settings passed to handlers via context.
type Config struct {
	Dialect        string
	Getenv         func(string) string
	dialectConfigs map[string]any
}

// DialectConfig returns opaque per-dialect configuration stored during option application.
func (c *Config) DialectConfig(name string) any {
	if c == nil || c.dialectConfigs == nil {
		return nil
	}
	return c.dialectConfigs[name]
}

// SetDialectConfig stores opaque per-dialect configuration.
func (c *Config) SetDialectConfig(name string, cfg any) {
	if c == nil {
		return
	}
	if c.dialectConfigs == nil {
		c.dialectConfigs = make(map[string]any)
	}
	c.dialectConfigs[name] = cfg
}

type configKey struct{}

// WithConfig attaches cfg to ctx for handler producers.
func WithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey{}, cfg)
}

// ConfigFrom returns the Config attached to ctx, or a default Config.
func ConfigFrom(ctx context.Context) *Config {
	if ctx == nil {
		return defaultConfig()
	}
	if cfg, ok := ctx.Value(configKey{}).(*Config); ok && cfg != nil {
		return cfg
	}
	return defaultConfig()
}

func defaultConfig() *Config {
	return &Config{Getenv: os.Getenv}
}

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
	cfg := defaultConfig()
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
