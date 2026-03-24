package core

import "strings"

// WatExecutionContext carries CLI-level settings for this wat invocation that
// hook handlers may consult alongside hook JSON (host, subcommand, optional run filters, etc.).
type WatExecutionContext struct {
	host       string
	subcommand string
	// filePattern is nil when no -f/--file-pattern was provided (shared argv, parsed with host after the subcommand).
	filePattern *string
}

// NewWatExecutionContext builds a context with the hook host name from argv (-H / --host).
func NewWatExecutionContext(host string) WatExecutionContext {
	return WatExecutionContext{host: host}
}

// Host returns the hook host (e.g. cursor).
func (c WatExecutionContext) Host() string {
	return c.host
}

// Subcommand returns the wat subcommand (e.g. run).
func (c WatExecutionContext) Subcommand() string {
	return c.subcommand
}

// FilePattern returns the optional regexp from shared argv (-f / --file-pattern), or nil if unset.
// When non-nil, the pointed-to string is already trimmed.
func (c WatExecutionContext) FilePattern() *string {
	return c.filePattern
}

// WithSubcommand returns a copy with the subcommand name set.
func (c WatExecutionContext) WithSubcommand(name string) WatExecutionContext {
	c.subcommand = name
	return c
}

// WithFilePattern returns a copy with the parsed regexp source. Empty or whitespace-only clears the pattern (nil).
func (c WatExecutionContext) WithFilePattern(pattern string) WatExecutionContext {
	s := strings.TrimSpace(pattern)
	if s == "" {
		c.filePattern = nil
		return c
	}
	c.filePattern = &s
	return c
}
