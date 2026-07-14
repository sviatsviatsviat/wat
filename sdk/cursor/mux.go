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

func registerObserveHandler[E Event](fn func(context.Context, Hook[E]) error) {
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
		if err := fn(ctx, NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, HandlerErrorExit, err
		}
		return nil, 0, nil
	})
}

func registerAny(fn func(context.Context, AnyHook) error) {
	if fn == nil {
		return
	}
	run.RegisterAnyHandler("cursor", "cursor", func(ctx context.Context, raw []byte) ([]byte, int, error) {
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

// PreToolUse registers a preToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, PreToolUseHook, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return &Chain{}
}

// PostToolUse registers a postToolUse handler.
func (c *Chain) PostToolUse(fn func(context.Context, PostToolUseHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUse) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// PostToolUseFailure registers a postToolUseFailure handler.
func (c *Chain) PostToolUseFailure(fn func(context.Context, PostToolUseFailureHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PostToolUseFailure) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// BeforeShellExecution registers a beforeShellExecution handler.
func (c *Chain) BeforeShellExecution(fn func(context.Context, BeforeShellExecutionHook, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeShellExecution) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return &Chain{}
}

// BeforeMCPExecution registers a beforeMCPExecution handler.
func (c *Chain) BeforeMCPExecution(fn func(context.Context, BeforeMCPExecutionHook, PermissionResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeMCPExecution) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionResults{})
	})
	return &Chain{}
}

// BeforeReadFile registers a beforeReadFile handler.
func (c *Chain) BeforeReadFile(fn func(context.Context, BeforeReadFileHook, BeforeReadFileResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeReadFile) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeReadFileResults{})
	})
	return &Chain{}
}

// BeforeTabFileRead registers a beforeTabFileRead handler.
func (c *Chain) BeforeTabFileRead(fn func(context.Context, BeforeTabFileReadHook, BeforeTabFileReadResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeTabFileRead) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeTabFileReadResults{})
	})
	return &Chain{}
}

// AfterShellExecution registers an afterShellExecution handler.
func (c *Chain) AfterShellExecution(fn func(context.Context, AfterShellExecutionHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterShellExecution) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// AfterMCPExecution registers an afterMCPExecution handler.
func (c *Chain) AfterMCPExecution(fn func(context.Context, AfterMCPExecutionHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterMCPExecution) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// AfterFileEdit registers an afterFileEdit handler.
func (c *Chain) AfterFileEdit(fn func(context.Context, AfterFileEditHook, PostToolResults) (PostToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev AfterFileEdit) (PostToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), postToolResults{})
	})
	return &Chain{}
}

// BeforeSubmitPrompt registers a beforeSubmitPrompt handler.
func (c *Chain) BeforeSubmitPrompt(fn func(context.Context, BeforeSubmitPromptHook, BeforeSubmitPromptResults) (BeforeSubmitPromptOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev BeforeSubmitPrompt) (BeforeSubmitPromptOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), beforeSubmitPromptResults{})
	})
	return &Chain{}
}

// Stop registers a stop handler.
func (c *Chain) Stop(fn func(context.Context, StopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}

// SubagentStop registers a subagentStop handler.
func (c *Chain) SubagentStop(fn func(context.Context, SubagentStopHook, StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}

// SubagentStart registers a subagentStart handler.
func (c *Chain) SubagentStart(fn func(context.Context, SubagentStartHook, SubagentStartResults) (PermissionOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (PermissionOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return &Chain{}
}

// SessionStart registers a sessionStart handler.
func (c *Chain) SessionStart(fn func(context.Context, SessionStartHook, SessionStartResults) (SessionStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SessionStart) (SessionStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), sessionStartResults{})
	})
	return &Chain{}
}

// PreCompact registers a preCompact handler.
func (c *Chain) PreCompact(fn func(context.Context, PreCompactHook, PreCompactResults) (PreCompactOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreCompact) (PreCompactOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preCompactResults{})
	})
	return &Chain{}
}

// SessionEnd registers an observe-only sessionEnd handler.
func (c *Chain) SessionEnd(fn func(context.Context, SessionEndHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// AfterAgentResponse registers an observe-only afterAgentResponse handler.
func (c *Chain) AfterAgentResponse(fn func(context.Context, AfterAgentResponseHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// AfterAgentThought registers an observe-only afterAgentThought handler.
func (c *Chain) AfterAgentThought(fn func(context.Context, AfterAgentThoughtHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// AfterTabFileEdit registers an observe-only afterTabFileEdit handler.
func (c *Chain) AfterTabFileEdit(fn func(context.Context, AfterTabFileEditHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// WorkspaceOpen registers an observe-only workspaceOpen handler.
func (c *Chain) WorkspaceOpen(fn func(context.Context, WorkspaceOpenHook) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}

// OnAny registers an observe-only handler invoked for every event.
func (c *Chain) OnAny(fn func(context.Context, AnyHook) error) *Chain {
	registerAny(fn)
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
