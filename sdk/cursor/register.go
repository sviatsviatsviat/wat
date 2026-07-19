package cursor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

var codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)

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

func registerHandler[E run.Event, O Output](fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event run.Event) ([]byte, int, error) {
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
		return encode(name, result)
	})
}

func registerObserveHandler[E run.Event](fn func(context.Context, run.Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	run.RegisterHandler(Dialect, name, func(ctx context.Context, event run.Event) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("cursor: handler for %s received %T", name, event)
		}
		if err := fn(ctx, run.NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, HandlerErrorExit, err
		}
		return nil, 0, nil
	})
}
