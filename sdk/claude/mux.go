package claude

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
	run.RegisterDialect("claude", run.DialectOps{
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
	if has("cursor_version") || has("conversation_id") || has("sessionId") {
		return false
	}
	if has("hook_event_name") && has("timestamp") {
		return false
	}
	return has("session_id")
}

func eventNameFromRaw(raw []byte, eventHint string) (string, error) {
	ev, err := Decode(raw)
	if err != nil {
		return "", err
	}
	name := ev.EventName()
	if name == "" {
		name = eventHint
	}
	if name == "" {
		return "", fmt.Errorf("claude: empty event name")
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
		panic(fmt.Sprintf("claude: duplicate handler for %s", name))
	}
	registered[name] = true
	registeredMu.Unlock()

	run.RegisterHandler("claude", "claude", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		ev, err := Decode(raw)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		typed, ok := ev.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, ev)
		}
		result, err := fn(ctx, typed)
		if err != nil {
			return nil, handlerErrorExit(ctx, name), err
		}
		if isZeroOutput(result) {
			return nil, 0, nil
		}
		cfg := claudeConfigFrom(ctx)
		stdout, err := Encode(name, result, cfg.encodeOpts()...)
		return stdout, 0, err
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
		panic(fmt.Sprintf("claude: duplicate handler for %s", name))
	}
	registered[name] = true
	registeredMu.Unlock()

	run.RegisterHandler("claude", "claude", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		ev, err := Decode(raw)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		typed, ok := ev.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, ev)
		}
		if err := fn(ctx, NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, handlerErrorExit(ctx, name), err
		}
		return nil, 0, nil
	})
}

func registerAny(fn func(context.Context, AnyHook) error) {
	if fn == nil {
		return
	}
	run.RegisterAnyHandler("claude", "claude", func(ctx context.Context, raw []byte) ([]byte, int, error) {
		ev, err := Decode(raw)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		if err := fn(ctx, NewAnyHook(run.InvocationFrom(ctx), ev)); err != nil {
			return nil, handlerErrorExit(ctx, ev.EventName()), err
		}
		return nil, 0, nil
	})
}

func handlerErrorExit(ctx context.Context, eventName string) int {
	cfg := claudeConfigFrom(ctx)
	if cfg.policy == FailBlock {
		return FailBlockExit
	}
	return HandlerErrorExit
}

func claudeConfigFrom(ctx context.Context) runtimeConfig {
	cfg := run.ConfigFrom(ctx)
	if v := cfg.DialectConfig("claude"); v != nil {
		if rc, ok := v.(*runtimeConfig); ok && rc != nil {
			return *rc
		}
	}
	return defaultRuntimeConfig()
}

// PreToolUse registers a PreToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, PreToolUseHook, PreToolUseResults) (PreToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PreToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preToolUseResults{})
	})
	return &Chain{}
}

// PostToolUse registers a PostToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolUseResults) (PostToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolUseResults{})
	})
	return &Chain{}
}

// PostToolUseFailure registers a PostToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailureHook, PostToolUseFailureResults) (PostToolUseOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolUseOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolUseFailureResults{})
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

// UserPromptSubmit registers a UserPromptSubmit handler.
func (c *Chain) UserPromptSubmit(fn func(context.Context, UserPromptSubmitHook, UserPromptSubmitResults) (UserPromptSubmitOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev UserPromptSubmit) (UserPromptSubmitOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), userPromptSubmitResults{})
	})
	return &Chain{}
}

// Stop registers a Stop handler.
func (c *Chain) Stop(fn func(context.Context, StopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
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
func (c *Chain) SubagentStart(fn func(context.Context, SubagentStartHook, SubagentStartResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return &Chain{}
}

// Notification registers a Notification handler.
func (c *Chain) Notification(fn func(context.Context, NotificationHook, NotificationResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Notification) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), notificationResults{})
	})
	return &Chain{}
}

// PreCompact registers a PreCompact handler.
func (c *Chain) PreCompact(fn func(context.Context, PreCompactHook, PreCompactResults) (CommonOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreCompact) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preCompactResults{})
	})
	return &Chain{}
}

// SessionEnd registers an observe-only SessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, SessionEndHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// OnAny registers an observe-only handler invoked for every event.
func (c *Chain) OnAny(fn func(context.Context, AnyHook) error) *Chain {
	registerAny(fn)
	return &Chain{}
}

// ResetHandlers clears claude registration tracking and claude-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	registeredMu.Lock()
	registered = make(map[string]bool)
	registeredMu.Unlock()
	run.ResetOwner("claude")
}
