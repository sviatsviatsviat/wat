package copilot

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
	run.RegisterDialect("copilot", run.DialectOps{
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
		return false
	}
	return has("sessionId") || (has("hook_event_name") && has("timestamp"))
}

func eventNameFromRaw(raw []byte, eventHint string) (string, error) {
	ev, err := decodeWithHint(raw, eventHint)
	if err != nil {
		return "", err
	}
	name := ev.EventName()
	if name == "" {
		return "", fmt.Errorf("copilot: empty event name")
	}
	return name, nil
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
		panic(fmt.Sprintf("copilot: duplicate handler for %s", name))
	}
	registered[name] = true
	registeredMu.Unlock()

	run.RegisterHandler("copilot", "copilot", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		ev, err := decodeWithHint(raw, cfg.EventHint)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		typed, ok := ev.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("copilot: handler for %s received %T", name, ev)
		}
		result, err := fn(ctx, typed)
		if err != nil {
			return nil, handlerErrorExit(name), err
		}
		if isZeroOutput(result) {
			return nil, 0, nil
		}
		return Encode(name, result)
	})
}

func registerObserveHandler[E Event](fn func(context.Context, Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	registeredMu.Lock()
	if registered[name] {
		registeredMu.Unlock()
		panic(fmt.Sprintf("copilot: duplicate handler for %s", name))
	}
	registered[name] = true
	registeredMu.Unlock()

	run.RegisterHandler("copilot", "copilot", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		ev, err := decodeWithHint(raw, cfg.EventHint)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		typed, ok := ev.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("copilot: handler for %s received %T", name, ev)
		}
		if err := fn(ctx, NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, handlerErrorExit(name), err
		}
		return nil, 0, nil
	})
}

func registerAny(fn func(context.Context, AnyHook) error) {
	if fn == nil {
		return
	}
	run.RegisterAnyHandler("copilot", "copilot", func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		ev, err := decodeWithHint(raw, cfg.EventHint)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		if err := fn(ctx, NewAnyHook(run.InvocationFrom(ctx), ev)); err != nil {
			return nil, HandlerErrorExit, err
		}
		return nil, 0, nil
	})
}

func handlerErrorExit(eventName string) int {
	if eventName == EventPreToolUse {
		return PreToolErrorExit
	}
	return HandlerErrorExit
}

// PreToolUse registers a PreToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, PreToolUseHook, PreToolResults) (PreToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PreToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preToolResults{})
	})
	return &Chain{}
}

// PostToolUse registers a PostToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// PostToolUseFailure registers a PostToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailureHook, PostToolFailureResults) (PostToolFailureOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolFailureOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolFailureResults{})
	})
	return &Chain{}
}

// AgentStop registers an AgentStop handler.
func (c *Chain) AgentStop(fn func(context.Context, AgentStopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AgentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}

// SubagentStop registers a SubagentStop handler.
func (c *Chain) SubagentStop(fn func(context.Context, SubagentStopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}

// PermissionRequest registers a PermissionRequest handler.
func (c *Chain) PermissionRequest(fn func(context.Context, PermissionRequestHook, PermissionRequestResults) (PermissionRequestOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionRequest) (PermissionRequestOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionRequestResults{})
	})
	return &Chain{}
}

// SessionStart registers a SessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStartHook, SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return &Chain{}
}

// SubagentStart registers a SubagentStart handler.
func (c *Chain) SubagentStart(fn func(context.Context, SubagentStartHook, SubagentStartResults) (SubagentStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (SubagentStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return &Chain{}
}

// Notification registers a Notification handler.
func (c *Chain) Notification(fn func(context.Context, NotificationHook, NotificationResults) (NotificationOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Notification) (NotificationOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), notificationResults{})
	})
	return &Chain{}
}

// SessionEnd registers an observe-only SessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, SessionEndHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// UserPromptSubmitted registers an observe-only userPromptSubmitted handler.
func (c *Chain) UserPromptSubmitted(fn func(context.Context, UserPromptSubmittedHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// PreCompact registers an observe-only preCompact handler.
func (c *Chain) PreCompact(fn func(context.Context, PreCompactHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// ErrorOccurred registers an observe-only errorOccurred handler.
func (c *Chain) ErrorOccurred(fn func(context.Context, ErrorOccurredHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// OnAny registers an observe-only handler invoked for every event.
func (c *Chain) OnAny(fn func(context.Context, AnyHook) error) *Chain {
	registerAny(fn)
	return &Chain{}
}

// ResetHandlers clears copilot registration tracking and copilot-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
	run.ResetOwner("copilot")
}
