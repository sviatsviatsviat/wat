package cursor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func init() {
	run.RegisterDialect(Dialect, run.DialectOps{
		Detect: detectPayload,
		Codec:  codec,
		Merge:  MergeOutputs,
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

func registerHandler[E Event, O any](fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event any) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("cursor: handler for %s received %T", name, event)
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

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event any) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("cursor: handler for %s received %T", name, event)
		}
		if err := fn(ctx, NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, HandlerErrorExit, err
		}
		return nil, 0, nil
	})
}
