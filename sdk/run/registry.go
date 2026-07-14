package run

import (
	"context"
	"sync"
)

// Producer handles one hook invocation: decode, run logic, and encode to native JSON.
type Producer func(ctx context.Context, raw []byte) (output []byte, exit int, err error)

// DialectOps supplies dialect-specific detection, event naming, and output merge.
type DialectOps struct {
	// Detect reports whether raw matches this dialect.
	Detect func(raw []byte, getenv func(string) string) bool
	// EventName returns the native event name for raw using eventHint when needed.
	EventName func(raw []byte, eventHint string) (string, error)
	// Merge combines native JSON outputs from multiple handlers for one event.
	Merge func(outputs [][]byte) ([]byte, error)
}

type ownedProducer struct {
	owner string
	p     Producer
}

// Registry holds registered dialects and handlers.
type Registry struct {
	mu sync.RWMutex

	dialectOrder []string
	dialects     map[string]DialectOps
	anyHandlers  map[string][]ownedProducer
	handlers     map[string]map[string][]ownedProducer
}

var defaultRegistry = NewRegistry()

// NewRegistry returns an empty handler registry.
func NewRegistry() *Registry {
	return &Registry{
		dialects:    make(map[string]DialectOps),
		anyHandlers: make(map[string][]ownedProducer),
		handlers:    make(map[string]map[string][]ownedProducer),
	}
}

// Default returns the package-level singleton registry.
func Default() *Registry {
	return defaultRegistry
}

// Reset clears handlers on the default registry but preserves registered dialects.
// It is intended for tests.
func Reset() {
	defaultRegistry.resetHandlers()
}

// ResetDialect clears handlers for dialect on the default registry but preserves
// dialect ops and handlers for other dialects. It is intended for tests.
func ResetDialect(dialect string) {
	defaultRegistry.resetDialectHandlers(dialect)
}

// ResetOwner removes handlers registered by owner on the default registry.
// It is intended for tests.
func ResetOwner(owner string) {
	defaultRegistry.resetOwner(owner)
}

// RegisterDialect registers dialect ops. Duplicate names panic.
func RegisterDialect(name string, ops DialectOps) {
	defaultRegistry.RegisterDialect(name, ops)
}

// RegisterHandler appends a handler for owner, dialect, and native event name.
func RegisterHandler(owner, dialect, eventName string, p Producer) {
	defaultRegistry.RegisterHandler(owner, dialect, eventName, p)
}

// RegisterAnyHandler appends a catch-all handler for owner on dialect.
func RegisterAnyHandler(owner, dialect string, p Producer) {
	defaultRegistry.RegisterAnyHandler(owner, dialect, p)
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
		r.handlers[name] = make(map[string][]ownedProducer)
	}
}

// RegisterHandler appends a handler for owner, dialect, and native event name on r.
func (r *Registry) RegisterHandler(owner, dialect, eventName string, p Producer) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.handlers[dialect] == nil {
		r.handlers[dialect] = make(map[string][]ownedProducer)
	}
	r.handlers[dialect][eventName] = append(r.handlers[dialect][eventName], ownedProducer{owner: owner, p: p})
}

// RegisterAnyHandler appends a catch-all handler for owner on dialect on r.
func (r *Registry) RegisterAnyHandler(owner, dialect string, p Producer) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.anyHandlers[dialect] = append(r.anyHandlers[dialect], ownedProducer{owner: owner, p: p})
}

func (r *Registry) producers(dialect, eventName string) []Producer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Producer
	for _, h := range r.anyHandlers[dialect] {
		out = append(out, h.p)
	}
	for _, h := range r.handlers[dialect][eventName] {
		out = append(out, h.p)
	}
	return out
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
	r.anyHandlers = make(map[string][]ownedProducer)
	for name := range r.handlers {
		r.handlers[name] = make(map[string][]ownedProducer)
	}
}

func (r *Registry) resetDialectHandlers(dialect string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.anyHandlers, dialect)
	if r.handlers[dialect] != nil {
		r.handlers[dialect] = make(map[string][]ownedProducer)
	}
}

func (r *Registry) resetOwner(owner string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for dialect, list := range r.anyHandlers {
		r.anyHandlers[dialect] = removeOwner(list, owner)
	}
	for dialect, events := range r.handlers {
		for event, list := range events {
			filtered := removeOwner(list, owner)
			if len(filtered) == 0 {
				delete(events, event)
			} else {
				events[event] = filtered
			}
		}
		if len(events) == 0 {
			r.handlers[dialect] = make(map[string][]ownedProducer)
		}
	}
}

func removeOwner(list []ownedProducer, owner string) []ownedProducer {
	if owner == "" {
		return list
	}
	out := list[:0]
	for _, h := range list {
		if h.owner != owner {
			out = append(out, h)
		}
	}
	return out
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
