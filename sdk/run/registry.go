package run

import (
	"sync"
)

// Codec peeks event names and decodes payloads for one agent dialect.
type Codec interface {
	// EventName returns the native event name for raw.
	// It should be a cheap discriminant peek, not a full typed decode.
	EventName(raw []byte) (string, error)
	// Decode parses raw into a dialect-specific event value for handler dispatch.
	Decode(raw []byte) (event Event, err error)
}

// DialectOps supplies dialect-specific detection, codec, and output merge.
type DialectOps struct {
	// Detect reports whether raw matches this dialect.
	Detect func(raw []byte, getenv func(string) string) bool
	// Codec peeks event names and decodes typed events for this dialect.
	Codec Codec
	// Merge combines native JSON outputs from multiple handlers for one event.
	Merge func(outputs [][]byte) ([]byte, error)
}

// Registry holds registered dialects and handlers.
type Registry struct {
	mu sync.RWMutex

	dialectOrder []string
	dialects     map[string]DialectOps
	handlers     map[string]map[string][]hookRegistration
}

var defaultRegistry = NewRegistry()

// NewRegistry returns an empty handler registry.
func NewRegistry() *Registry {
	return &Registry{
		dialects: make(map[string]DialectOps),
		handlers: make(map[string]map[string][]hookRegistration),
	}
}

// GetDefaultRegistry returns the process-wide registry used by Main when
// callers register via UseHooks with no registry argument.
func GetDefaultRegistry() *Registry {
	return defaultRegistry
}

// registerDialect registers dialect ops on r. Duplicate names panic.
func (r *Registry) registerDialect(name string, ops DialectOps) {
	if name == "" {
		panic("run: registerDialect: empty name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dialects[name]; exists {
		panic("run: duplicate dialect " + name)
	}
	r.dialectOrder = append(r.dialectOrder, name)
	r.dialects[name] = ops
	if r.handlers[name] == nil {
		r.handlers[name] = make(map[string][]hookRegistration)
	}
}

// EnsureDialect registers dialect ops on r if name is not already present.
func (r *Registry) EnsureDialect(name string, ops DialectOps) {
	if name == "" {
		panic("run: EnsureDialect: empty name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dialects[name]; exists {
		return
	}
	r.dialectOrder = append(r.dialectOrder, name)
	r.dialects[name] = ops
	if r.handlers[name] == nil {
		r.handlers[name] = make(map[string][]hookRegistration)
	}
}

// RegisterHandler appends a typed output-producing hook registration for dialect on r.
func (r *Registry) RegisterHandler(dialect string, reg hookRegistration) {
	r.register(dialect, reg)
}

// RegisterObserveHandler appends a typed observe-only hook registration for dialect on r.
func (r *Registry) RegisterObserveHandler(dialect string, reg hookRegistration) {
	r.register(dialect, reg)
}

func (r *Registry) register(dialect string, reg hookRegistration) {
	if reg == nil {
		return
	}
	eventName := reg.eventName()
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.handlers[dialect] == nil {
		r.handlers[dialect] = make(map[string][]hookRegistration)
	}
	r.handlers[dialect][eventName] = append(r.handlers[dialect][eventName], reg)
}

func (r *Registry) handlersFor(dialect, eventName string) []hookRegistration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]hookRegistration(nil), r.handlers[dialect][eventName]...)
}

func (r *Registry) dialectOps(name string) (DialectOps, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ops, ok := r.dialects[name]
	return ops, ok
}

func (r *Registry) detectDialect(raw []byte, getenv func(string) string, forced string) (string, bool) {
	if forced != "" {
		if _, ok := r.dialectOps(forced); ok {
			return forced, true
		}
		return "", false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, name := range r.dialectOrder {
		ops := r.dialects[name]
		if ops.Detect != nil && ops.Detect(raw, getenv) {
			return name, true
		}
	}
	return "", false
}
