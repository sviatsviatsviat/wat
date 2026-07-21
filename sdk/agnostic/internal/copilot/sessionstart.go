package copilot

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// RegisterSessionStart registers fn on the Copilot SessionStart chain.
func RegisterSessionStart(registry *run.Registry, fn model.SessionStartHandler) {
	if fn == nil {
		return
	}
	sdkcopilot.UseHooks(registry).SessionStart(func(ctx context.Context, hook run.Hook[sdkcopilot.SessionStart], native sdkcopilot.SessionStartResults) (sdkcopilot.SessionStartOutput, error) {
		out, err := fn(ctx, model.NewSessionStartHook(hook.Invocation(), mapSessionStart(hook.Event)), newSessionStartResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapSessionStart(out)
		if !ok {
			return nil, fmt.Errorf("copilot: SessionStart result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapSessionStart(e sdkcopilot.SessionStart) *model.SessionStartEvent {
	return &model.SessionStartEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Life:     &model.Lifecycle{Source: e.Source, InitialPrompt: e.InitialPrompt()},
	}
}

func newSessionStartResults(native sdkcopilot.SessionStartResults) model.SessionStartResults {
	return sessionStartResults{native: native}
}

func unwrapSessionStart(r model.SessionStartResult) (sdkcopilot.SessionStartOutput, bool) {
	out, ok := r.(sessionStartResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type sessionStartResults struct {
	native sdkcopilot.SessionStartResults
}

// Context returns a context-injection result.
func (w sessionStartResults) Context(text string) model.SessionStartResult {
	return sessionStartResult{native: w.native.Context(text)}
}

type sessionStartResult struct {
	native sdkcopilot.SessionStartOutput
}

// IsZero reports whether the result carries no instruction.
func (r sessionStartResult) IsZero() bool { return r.native.IsZero() }
