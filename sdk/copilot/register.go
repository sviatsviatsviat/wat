package copilot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func init() {
	run.RegisterDialect(Dialect, run.DialectOps{
		Detect:    detectPayload,
		EventName: eventNameFromRaw,
		Decode: func(raw []byte, hint string) (any, error) {
			return decodeWithHint(raw, hint)
		},
		Merge: MergeOutputs,
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

func registerHandler[E Event, O any](fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event any) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("copilot: handler for %s received %T", name, event)
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

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event any) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("copilot: handler for %s received %T", name, event)
		}
		if err := fn(ctx, NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, handlerErrorExit(name), err
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
