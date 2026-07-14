package copilot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/copilot/internal"
	"github.com/sviatsviatsviat/wat/sdk/run"
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

	internal.MarkRegistered("copilot", name)

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

	internal.MarkRegistered("copilot", name)

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

// ResetHandlers clears registration tracking and copilot-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	internal.ResetRegistered()
	run.ResetOwner("copilot")
}
