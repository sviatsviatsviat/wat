package claude

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// RegisterSessionStart registers fn on the Claude SessionStart chain.
func RegisterSessionStart(fn model.SessionStartHandler) {
	if fn == nil {
		return
	}
	sdkclaude.UseHooks().SessionStart(func(ctx context.Context, hook sdkclaude.SessionStart, native sdkclaude.SessionStartResults) (sdkclaude.SessionStartOutput, error) {
		out, err := fn(ctx, *mapSessionStart(hook), newSessionStartResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapSessionStart(out)
		if !ok {
			return nil, fmt.Errorf("claude: SessionStart result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapSessionStart(e sdkclaude.SessionStart) *model.SessionStartEvent {
	return &model.SessionStartEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Life:     &model.Lifecycle{Source: e.Source, Model: e.Model},
	}
}

func newSessionStartResults(native sdkclaude.SessionStartResults) model.SessionStartResults {
	return sessionStartResults{native: native}
}

func unwrapSessionStart(r model.SessionStartResult) (sdkclaude.SessionStartOutput, bool) {
	out, ok := r.(sessionStartResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type sessionStartResults struct {
	native sdkclaude.SessionStartResults
}

// Context returns a context-injection result.
func (w sessionStartResults) Context(text string) model.SessionStartResult {
	return sessionStartResult{native: w.native.Context(text)}
}

type sessionStartResult struct {
	native sdkclaude.SessionStartOutput
}

// IsZero reports whether the result carries no instruction.
func (r sessionStartResult) IsZero() bool { return r.native.IsZero() }
