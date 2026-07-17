package internal

import (
	"fmt"
	"strings"
	"sync"
)

var (
	registeredMu sync.Mutex
	registered   = make(map[string]bool)
)

func registrationKey(owner, name string) string {
	return owner + "\x00" + name
}

// MarkRegistered records a handler registration or panics on duplicate for the same owner and event.
// Owner "agnostic" is exempt so portable On* can register multiple handlers that run merges.
func MarkRegistered(owner, name string) {
	if owner == "agnostic" {
		return
	}
	registeredMu.Lock()
	defer registeredMu.Unlock()
	key := registrationKey(owner, name)
	if registered[key] {
		panic(fmt.Sprintf("%s: duplicate handler for %s", owner, name))
	}
	registered[key] = true
}

// ResetRegistered clears all registration tracking.
func ResetRegistered() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
}

// ResetRegisteredOwner clears registration tracking for a single run owner.
func ResetRegisteredOwner(owner string) {
	registeredMu.Lock()
	defer registeredMu.Unlock()
	prefix := owner + "\x00"
	for k := range registered {
		if strings.HasPrefix(k, prefix) {
			delete(registered, k)
		}
	}
}
