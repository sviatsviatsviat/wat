package run

import (
	"context"
	"sync"
)

// Producer handles one hook invocation: run logic on an already-decoded event
// and encode native JSON. Serve decodes the payload once before calling producers.
type Producer func(ctx context.Context, event any) (output []byte, exit int, err error)

// DialectOps supplies dialect-specific detection, event naming, decode, and output merge.
type DialectOps struct {
	// Detect reports whether raw matches this dialect.
	Detect func(raw []byte, getenv func(string) string) bool
	// EventName returns the native event name for raw using eventHint when needed.
	// It should be a cheap discriminant peek, not a full typed decode.
	EventName func(raw []byte, eventHint string) (string, error)
	// Decode parses raw into a dialect-specific event value for handler dispatch.
	Decode func(raw []byte, eventHint string) (event any, err error)
	// Merge combines native JSON outputs from multiple handlers for one event.
	Merge func(outputs [][]byte) ([]byte, error)
}

// Registry holds registered dialects and handlers.
type Registry struct {
	mu sync.RWMutex

	dialectOrder []string
	dialects     map[string]DialectOps
	handlers     map[string]map[string][]Producer
}

var defaultRegistry = NewRegistry()

// NewRegistry returns an empty handler registry.
func NewRegistry() *Registry {
	return &Registry{
		dialects: make(map[string]DialectOps),
		handlers: make(map[string]map[string][]Producer),
	}
}

// Reset clears handlers on the default registry but preserves registered dialects.
// It is intended for tests.
func Reset() {
	defaultRegistry.resetHandlers()
}

// RegisterDialect registers dialect ops. Duplicate names panic.
func RegisterDialect(name string, ops DialectOps) {
	defaultRegistry.RegisterDialect(name, ops)
}

// RegisterHandler appends a handler for dialect and native event name.
func RegisterHandler(dialect, eventName string, p Producer) {
	defaultRegistry.RegisterHandler(dialect, eventName, p)
}

// RegisterDialect registers dialect ops on r. Duplicate names panic.
func (r *Registry) RegisterDialect(name string, ops DialectOps) {
	if name == "" {
		panic("run: RegisterDialect: empty name")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dialects[name]; exists {
		panic("run: duplicate dialect " + name)
	}
	r.dialectOrder = append(r.dialectOrder, name)
	r.dialects[name] = ops
	if r.handlers[name] == nil {
		r.handlers[name] = make(map[string][]Producer)
	}
}

// RegisterHandler appends a handler for dialect and native event name on r.
func (r *Registry) RegisterHandler(dialect, eventName string, p Producer) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.handlers[dialect] == nil {
		r.handlers[dialect] = make(map[string][]Producer)
	}
	r.handlers[dialect][eventName] = append(r.handlers[dialect][eventName], p)
}

func (r *Registry) producers(dialect, eventName string) []Producer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Producer(nil), r.handlers[dialect][eventName]...)
}

func (r *Registry) dialectOps(name string) (DialectOps, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ops, ok := r.dialects[name]
	return ops, ok
}

func (r *Registry) resetHandlers() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name := range r.handlers {
		r.handlers[name] = make(map[string][]Producer)
	}
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
