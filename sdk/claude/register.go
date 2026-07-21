package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

var codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired, newEncoder())

var (
	defaultReg   = run.GetDefaultRegistry()
	defaultChain = newChain(defaultReg)
)

func init() {
	ensureDialect(defaultReg)
}

func dialectOps() run.DialectOps {
	return run.DialectOps{
		Detect: detectPayload,
		Codec:  codec,
		Merge:  MergeOutputs,
	}
}

func ensureDialect(r *run.Registry) {
	r.EnsureDialect(Dialect, dialectOps())
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
	if has("hook_event_name") && has("timestamp") {
		return false
	}
	return has("session_id")
}

func registerHandler[E run.Event, O Output](r *run.Registry, fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	r.RegisterHandler(Dialect, name, func(ctx context.Context, event run.Event) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, event)
		}
		result, err := fn(ctx, typed)
		if err != nil {
			return nil, HandlerErrorExit, err
		}
		return codec.Encode(name, Output(result))
	})
}

func registerObserveHandler[E run.Event](r *run.Registry, fn func(context.Context, run.Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	r.RegisterHandler(Dialect, name, func(ctx context.Context, event run.Event) ([]byte, int, error) {
		typed, ok := event.(E)
		if !ok {
			return nil, HandlerErrorExit, fmt.Errorf("claude: handler for %s received %T", name, event)
		}
		if err := fn(ctx, run.NewHook(run.InvocationFrom(ctx), typed)); err != nil {
			return nil, HandlerErrorExit, err
		}
		return nil, SuccessExit, nil
	})
}
