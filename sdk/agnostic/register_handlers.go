package agnostic

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

type wireResult interface {
	Result() model.Result
}

func registerResultHandler[R wireResult](kind model.Kind, fn func(context.Context, *model.Event) (R, error)) {
	if fn == nil {
		return
	}
	wrap := func(ctx context.Context, ev *model.Event) (model.Result, error) {
		res, err := fn(ctx, ev)
		if err != nil {
			return model.Result{}, err
		}
		return res.Result(), nil
	}
	for _, agent := range []model.Dialect{model.Claude, model.Copilot, model.Cursor} {
		for _, eventName := range eventsForKind(agent, kind) {
			registerAgnosticHandler(agent, eventName, kind, wrap)
		}
	}
}

func registerObserveHandler(kind model.Kind, fn func(context.Context, *model.Event) error) {
	if fn == nil {
		return
	}
	wrap := func(ctx context.Context, ev *model.Event) (model.Result, error) {
		if err := fn(ctx, ev); err != nil {
			return model.Result{}, err
		}
		return model.Result{}, nil
	}
	for _, agent := range []model.Dialect{model.Claude, model.Copilot, model.Cursor} {
		for _, eventName := range eventsForKind(agent, kind) {
			registerAgnosticHandler(agent, eventName, kind, wrap)
		}
	}
}

func registerAny(fn AnyHandler) {
	if fn == nil {
		return
	}
	for _, agent := range []model.Dialect{model.Claude, model.Copilot, model.Cursor} {
		run.RegisterAnyHandler("agnostic", agent.String(), makeAnyProducer(agent, fn))
	}
}

func registerAgnosticHandler(agent model.Dialect, eventName string, kind model.Kind, fn func(context.Context, *model.Event) (model.Result, error)) {
	run.RegisterHandler("agnostic", agent.String(), eventName, func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, err
		}
		if ev.Kind != kind {
			return nil, 0, nil
		}
		res, err := fn(ctx, ev)
		if err != nil {
			return nil, handlerErrorExit(ev), err
		}
		stdout, code, err := codec.Encode(ev, res)
		return stdout, code, err
	})
}

func makeAnyProducer(agent model.Dialect, fn AnyHandler) run.Producer {
	return func(ctx context.Context, raw []byte) ([]byte, int, error) {
		cfg := run.ConfigFrom(ctx)
		codec, err := codecForServe(agent, cfg)
		if err != nil {
			return nil, 1, err
		}
		ev, err := codec.Decode(raw, cfg.EventHint)
		if err != nil {
			return nil, 1, err
		}
		if err := fn(ctx, anyHook(run.InvocationFrom(ctx), AnyEventFrom(ev))); err != nil {
			return nil, handlerErrorExit(ev), err
		}
		return nil, 0, nil
	}
}

func codecForServe(agent model.Dialect, _ *run.Config) (model.Codec, error) {
	switch agent {
	case model.Claude:
		return &claude.Codec{}, nil
	case model.Copilot:
		return &copilot.Codec{}, nil
	case model.Cursor:
		return &cursor.Codec{}, nil
	default:
		return nil, fmt.Errorf("agnostic: no codec for dialect %q", agent)
	}
}

func eventsForKind(agent model.Dialect, kind model.Kind) []string {
	switch agent {
	case model.Claude:
		if name, ok := claude.EventForKind[kind]; ok {
			return []string{name}
		}
	case model.Copilot:
		if name, ok := copilot.EventForKind[kind]; ok {
			return []string{name}
		}
	case model.Cursor:
		var out []string
		for event, k := range cursor.KindForEventMap {
			if k == kind {
				out = append(out, event)
			}
		}
		return out
	}
	return nil
}

func handlerErrorExit(ev *model.Event) int {
	switch ev.Agent {
	case model.Copilot:
		if ev.Kind == model.KindPreTool {
			return copilot.PreToolErrorExit
		}
		return 1
	case model.Cursor:
		return cursor.HandlerErrorExit
	default:
		return 1
	}
}
