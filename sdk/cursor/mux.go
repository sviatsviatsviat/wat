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
func (c *Chain) BeforeShellExecution(fn func(context.Context, BeforeShellExecution, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeShellExecution) (PermissionOutput, error) {
		return fn(ctx, ev, permissionResults{})
	})
	return &Chain{}
}

// BeforeMCPExecution registers a beforeMCPExecution handler.
func (c *Chain) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecution, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeMCPExecution) (PermissionOutput, error) {
		return fn(ctx, ev, permissionResults{})
	})
	return &Chain{}
}

// AfterFileEdit registers an afterFileEdit handler.
func (c *Chain) AfterFileEdit(fn func(context.Context, AfterFileEdit, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterFileEdit) (PostToolOutput, error) {
		return fn(ctx, ev, postToolResults{})
	})
	return &Chain{}
}

// BeforeSubmitPrompt registers a beforeSubmitPrompt handler.
func (c *Chain) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPrompt, BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeSubmitPrompt) (BeforeSubmitPromptOutput, error) {
		return fn(ctx, ev, beforeSubmitPromptResults{})
	})
	return &Chain{}
}

// Stop registers a stop handler.
func (c *Chain) Stop(fn func(context.Context, Stop, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, ev, stopResults{})
	})
	return &Chain{}
}

// SessionStart registers a sessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStart, SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, ev, sessionStartResults{})
	})
	return &Chain{}
}

// ResetHandlers clears cursor registration tracking and cursor-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
	run.ResetOwner("cursor")
}
