package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/claude/internal"
	"github.com/sviatsviatsviat/wat/sdk/run"
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

func registerHandler[E Event, O any](owner string, fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	internal.MarkRegistered(owner, name)

	run.RegisterHandler(owner, "claude", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
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
		rc := claudeRunConfig(run.ConfigFrom(ctx))
		stdout, err := Encode(name, result, rc.encodeOpts()...)
		return stdout, 0, err
	})
}

func registerObserveHandler[E Event](owner string, fn func(context.Context, Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	internal.MarkRegistered(owner, name)

	run.RegisterHandler(owner, "claude", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
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

func handlerErrorExit(ctx context.Context, _ string) int {
	rc := claudeRunConfig(run.ConfigFrom(ctx))
	if rc.policy == FailBlock {
		return FailBlockExit
	}
	return HandlerErrorExit
}

// ResetHandlers clears registration tracking and claude-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	internal.ResetRegisteredOwner("claude")
	run.ResetOwner("claude")
}

// ResetAdapter clears agnostic-owned registration tracking for this SDK.
// Pair with run.ResetOwner("agnostic") from sdk/agnostic.ResetHandlers.
func ResetAdapter() {
	internal.ResetRegisteredOwner(adapterOwner)
}
