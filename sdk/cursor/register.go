package cursor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal"
	"github.com/sviatsviatsviat/wat/sdk/run"
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

func registerHandler[E Event, O any](owner string, fn func(context.Context, E) (O, error)) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	internal.MarkRegistered(owner, name)

	run.RegisterHandler(owner, "cursor", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
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

func registerObserveHandler[E Event](owner string, fn func(context.Context, Hook[E]) error) {
	if fn == nil {
		return
	}
	var zero E
	name := zero.EventName()

	internal.MarkRegistered(owner, name)

	run.RegisterHandler(owner, "cursor", name, func(ctx context.Context, raw []byte) ([]byte, int, error) {
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

func registerAny(owner string, fn func(context.Context, AnyHook) error) {
	if fn == nil {
		return
	}
	run.RegisterAnyHandler(owner, "cursor", func(ctx context.Context, raw []byte) ([]byte, int, error) {
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

// ResetHandlers clears registration tracking and cursor-owned handlers
// in the shared run registry. It is intended for tests.
func ResetHandlers() {
	internal.ResetRegisteredOwner("cursor")
	run.ResetOwner("cursor")
}

// ResetAdapter clears agnostic-owned registration tracking for this SDK.
// Pair with run.ResetOwner("agnostic") from sdk/agnostic.ResetHandlers.
func ResetAdapter() {
	internal.ResetRegisteredOwner(adapterOwner)
}
