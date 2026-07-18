package cursor

import (
	"context"
	"fmt"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterSessionStart registers fn on the Cursor SessionStart chain.
func RegisterSessionStart(fn model.SessionStartHandler) {
	if fn == nil {
		return
	}
	sdkcursor.OnSessionStart(func(ctx context.Context, hook sdkcursor.Hook[sdkcursor.SessionStart], native sdkcursor.SessionStartResults) (sdkcursor.SessionStartOutput, error) {
		out, err := fn(ctx, model.NewSessionStartHook(hook.Invocation(), mapSessionStart(hook.Event, hook.Raw())), newSessionStartResults(native))
		if err != nil || out == nil {
			return nil, err
		}
		nativeOut, ok := unwrapSessionStart(out)
		if !ok {
			return nil, fmt.Errorf("cursor: SessionStart result must come from the injected Results builder")
		}
		return nativeOut, nil
	})
}

func mapSessionStart(e sdkcursor.SessionStart, raw []byte) *model.SessionStartEvent {
	return &model.SessionStartEvent{
		Envelope: envelope(e, raw),
		Life:     &model.Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent},
	}
}

func newSessionStartResults(native sdkcursor.SessionStartResults) model.SessionStartResults {
	return sessionStartResults{native: native}
}

func unwrapSessionStart(r model.SessionStartResult) (sdkcursor.SessionStartOutput, bool) {
	out, ok := r.(sessionStartResult)
	if !ok {
		return nil, false
	}
	return out.native, true
}

type sessionStartResults struct {
	native sdkcursor.SessionStartResults
}

// Context returns a context-injection result.
func (w sessionStartResults) Context(text string) model.SessionStartResult {
	return sessionStartResult{native: w.native.Context(text)}
}

type sessionStartResult struct {
	native sdkcursor.SessionStartOutput
}

// IsZero reports whether the result carries no instruction.
func (r sessionStartResult) IsZero() bool { return sdkcursor.IsZeroOutput(r.native) }
