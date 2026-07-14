package cursor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Chain supports fluent handler registration into the shared run registry.
type Chain struct{}

var (
	registeredMu sync.Mutex
	registered   = make(map[string]bool)
)

func init() {
	run.RegisterDialect("cursor", run.DialectOps{
		Detect:    detectPayload,
		EventName: eventNameFromRaw,
		Merge:     MergeOutputs,
	})
}

func detectPayload(raw []byte, getenv func(string) string) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(raw, &probe); err != nil {
		return false
	}
	has := func(k string) bool { _, ok := probe[k]; return ok }
	if has("cursor_version") || has("conversation_id") {
		return true
	}
	if getenv != nil && getenv("CURSOR_VERSION") != "" {
		return true
	}
	return false
}

func eventNameFromRaw(raw []byte, eventHint string) (string, error) {
	ev, err := decodeWithHint(raw, eventHint)
	if err != nil {
		return "", err
	}
	name := ev.EventName()
	if name == "" {
		return "", fmt.Errorf("cursor: empty event name")
	}
	return name, nil
}

// On registers a typed handler for event type E into the shared run registry.
func On[E Event, O any](fn func(context.Context, E) (O, error)) *Chain {
	registerHandler(fn)
	return &Chain{}
}

func registerHandler[E Event, O any](fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	registeredMu.Lock()
	if registered[name] {
		registeredMu.Unlock()
		panic(fmt.Sprintf("cursor: duplicate handler for %s", name))
	}
	registered[name] = true
	registeredMu.Unlock()

	run.RegisterHandler("cursor", "cursor", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		ev, err := decodeWithHint(raw, cfg.EventHint)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		typed, ok := ev.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("cursor: handler for %s received %T", name, ev)
		}
		result, err := fn(ctx, typed)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		if isZeroOutput(result) {
			return nil, 0, nil
		}
		return Encode(name, result)
	})
}

// BeforeShellExecution registers a beforeShellExecution handler.
func (c *Chain) BeforeShellExecution(fn func(context.Context, BeforeShellExecution) (PermissionOutput, error)) *Chain {
	return On(fn)
}

// BeforeMCPExecution registers a beforeMCPExecution handler.
func (c *Chain) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecution) (PermissionOutput, error)) *Chain {
	return On(fn)
}

// AfterFileEdit registers an afterFileEdit handler.
func (c *Chain) AfterFileEdit(fn func(context.Context, AfterFileEdit) (PostToolOutput, error)) *Chain {
	return On(fn)
}

// BeforeSubmitPrompt registers a beforeSubmitPrompt handler.
func (c *Chain) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPrompt) (BeforeSubmitPromptOutput, error)) *Chain {
	return On(fn)
}

// Stop registers a stop handler.
func (c *Chain) Stop(fn func(context.Context, Stop) (StopOutput, error)) *Chain {
	return On(fn)
}

// SessionStart registers a sessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStart) (SessionStartOutput, error)) *Chain {
	return On(fn)
}

// ResetHandlers clears cursor registration tracking and cursor-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
	run.ResetOwner("cursor")
}
