package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// UserPromptEvent is the normalized view of a UserPrompt hook invocation.
type UserPromptEvent struct {
	Envelope
	Prompt string
}

// UserPromptEventFrom maps a decoded Event to UserPromptEvent.
func UserPromptEventFrom(ev *model.Event) (UserPromptEvent, error) {
	if ev == nil {
		return UserPromptEvent{}, fmt.Errorf("agnostic: nil event")
	}
	if ev.Kind != model.KindUserPrompt {
		return UserPromptEvent{}, fmt.Errorf("agnostic: expected UserPrompt kind, got %s", ev.Kind)
	}
	return UserPromptEvent{Envelope: envelopeFrom(ev), Prompt: ev.Prompt}, nil
}

// UserPromptHook is the handler context for portable UserPrompt events.
type UserPromptHook struct {
	UserPromptEvent
	inv run.Invocation
}

// Invocation returns serve-time settings for this hook invocation.
func (h UserPromptHook) Invocation() run.Invocation { return h.inv }

// Raw returns the untouched native JSON payload.
func (h UserPromptHook) Raw() json.RawMessage { return h.UserPromptEvent.Raw }

func userPromptHook(ctx run.Invocation, ev UserPromptEvent) UserPromptHook {
	return UserPromptHook{UserPromptEvent: ev, inv: ctx}
}

// UserPromptHandler handles observe-only UserPrompt events.
type UserPromptHandler func(ctx context.Context, hook UserPromptHook) error

// OnUserPrompt registers an observe-only handler for UserPrompt events.
func OnUserPrompt(fn UserPromptHandler) *Chain {
	if fn == nil {
		return &Chain{}
	}
	registerObserveHandler(model.KindUserPrompt, func(ctx context.Context, ev *model.Event) error {
		typed, err := UserPromptEventFrom(ev)
		if err != nil {
			return err
		}
		return fn(ctx, userPromptHook(run.InvocationFrom(ctx), typed))
	})
	return &Chain{}
}

// OnUserPrompt registers another observe-only UserPrompt handler on the chain.
func (c *Chain) OnUserPrompt(fn UserPromptHandler) *Chain {
	return OnUserPrompt(fn)
}
