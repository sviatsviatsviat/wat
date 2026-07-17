package agnostic

import (
	"context"
	"encoding/json"
	"fmt"

	agclaude "github.com/sviatsviatsviat/wat/sdk/agnostic/claude"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
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
	sdkclaude.Adapter().UserPromptSubmit(adaptClaudeUserPrompt(fn))
	sdkcopilot.Adapter().UserPromptSubmitted(adaptCopilotUserPrompt(fn))
	sdkcursor.Adapter().BeforeSubmitPrompt(adaptCursorUserPrompt(fn))
	return &Chain{}
}

// OnUserPrompt registers another observe-only UserPrompt handler on the chain.
func (c *Chain) OnUserPrompt(fn UserPromptHandler) *Chain {
	return OnUserPrompt(fn)
}

func adaptClaudeUserPrompt(fn UserPromptHandler) func(context.Context, sdkclaude.Hook[sdkclaude.UserPromptSubmit], sdkclaude.UserPromptSubmitResults) (sdkclaude.UserPromptSubmitOutput, error) {
	return func(ctx context.Context, hook sdkclaude.Hook[sdkclaude.UserPromptSubmit], _ sdkclaude.UserPromptSubmitResults) (sdkclaude.UserPromptSubmitOutput, error) {
		typed, err := UserPromptEventFrom(agclaude.MapEvent(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, userPromptHook(hook.Invocation(), typed))
	}
}

func adaptCopilotUserPrompt(fn UserPromptHandler) func(context.Context, sdkcopilot.Hook[sdkcopilot.UserPromptSubmitted]) error {
	return func(ctx context.Context, hook sdkcopilot.Hook[sdkcopilot.UserPromptSubmitted]) error {
		typed, err := UserPromptEventFrom(agcopilot.MapEvent(hook.Event, hook.Raw()))
		if err != nil {
			return err
		}
		return fn(ctx, userPromptHook(hook.Invocation(), typed))
	}
}

func adaptCursorUserPrompt(fn UserPromptHandler) func(context.Context, sdkcursor.Hook[sdkcursor.BeforeSubmitPrompt], sdkcursor.BeforeSubmitPromptResults) (sdkcursor.BeforeSubmitPromptOutput, error) {
	return func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.BeforeSubmitPrompt], _ sdkcursor.BeforeSubmitPromptResults) (sdkcursor.BeforeSubmitPromptOutput, error) {
		typed, err := UserPromptEventFrom(agcursor.MapEvent(hook.Event, hook.Raw()))
		if err != nil {
			return nil, err
		}
		return nil, fn(ctx, userPromptHook(hook.Invocation(), typed))
	}
}
