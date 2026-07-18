package claude

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
		Decode: func(raw []byte, _ string) (any, error) {
			return decode(raw)
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
	if has("cursor_version") || has("conversation_id") || has("sessionId") {
		return false
	}
	if has("hook_event_name") && has("timestamp") {
		return false
	}
	return has("session_id")
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
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, event)
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

func registerObserveHandler[E Event](fn func(context.Context, Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event any) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, event)
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
