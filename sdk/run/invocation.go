package run

import "context"

// Invocation exposes serve-time settings to hook handlers.
type Invocation struct {
	cfg *Config
}

// InvocationFrom returns serve-time settings attached to ctx.
func InvocationFrom(ctx context.Context) Invocation {
	return Invocation{cfg: ConfigFrom(ctx)}
}

// Dialect returns the forced dialect name, if any.
func (i Invocation) Dialect() string {
	if i.cfg == nil {
		return ""
	}
	return i.cfg.Dialect
}

// EventHint returns the native event name hint for payloads that omit it.
func (i Invocation) EventHint() string {
	if i.cfg == nil {
		return ""
	}
	return i.cfg.EventHint
}

// Getenv looks up an environment variable using the serve-time getenv function.
func (i Invocation) Getenv(key string) string {
	if i.cfg == nil || i.cfg.Getenv == nil {
		return ""
	}
	return i.cfg.Getenv(key)
}

// DialectConfig returns opaque per-dialect configuration stored during option application.
func (i Invocation) DialectConfig(name string) any {
	if i.cfg == nil {
		return nil
	}
	return i.cfg.DialectConfig(name)
}
