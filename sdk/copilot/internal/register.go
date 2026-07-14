package internal

import (
	"fmt"
	"sync"
)

var (
	registeredMu sync.Mutex
	registered   = make(map[string]bool)
)

// MarkRegistered records a handler registration or panics on duplicate.
func MarkRegistered(agent, name string) {
	registeredMu.Lock()
	defer registeredMu.Unlock()
	if registered[name] {
		panic(fmt.Sprintf("%s: duplicate handler for %s", agent, name))
	}
	registered[name] = true
}

// ResetRegistered clears registration tracking.
func ResetRegistered() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
}
