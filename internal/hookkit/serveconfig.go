package hookkit

import (
	"context"
	"os"
)

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

// WithConfig attaches cfg to ctx for hook handlers.
func WithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey{}, cfg)
}

// ConfigFrom returns the Config attached to ctx, or a default Config.
func ConfigFrom(ctx context.Context) *Config {
	if ctx == nil {
		return DefaultConfig()
	}
	if cfg, ok := ctx.Value(configKey{}).(*Config); ok && cfg != nil {
		return cfg
	}
	return DefaultConfig()
}

// DefaultConfig returns a Config with os.Getenv as the environment lookup.
func DefaultConfig() *Config {
	return &Config{Getenv: os.Getenv}
}
