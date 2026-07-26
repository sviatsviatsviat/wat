package cursor

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"

	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// RegisterPostToolFailure registers fn on the Cursor PostToolUseFailure chain.
//
// Cursor postToolUseFailure is observe-only. Portable Context results are
// accepted by the handler signature but discarded because the host documents
// no output fields for this event.
func RegisterPostToolFailure(fn model.PostToolFailureHandler) run.Hooks {
	if fn == nil {
		return nil
	}
	return sdkcursor.UseHooks().PostToolUseFailure(func(ctx context.Context, hook sdkcursor.PostToolUseFailure) error {
		_, err := fn(ctx, *mapPostToolUseFailure(hook), ignoredPostToolFailureResults{})
		return err
	})
}

func mapPostToolUseFailure(e sdkcursor.PostToolUseFailure) *model.PostToolFailureEvent {
	return &model.PostToolFailureEvent{
		Envelope: envelope(e.Envelope, e.EventName()),
		Tool:     model.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID),
		Result: &model.ToolResult{
			Error:       e.ErrorMessage,
			FailureType: e.FailureType,
			DurationMs:  e.DurationMillis(),
			IsInterrupt: e.IsInterrupt,
		},
	}
}

// ignoredPostToolFailureResults satisfies the portable builder while dropping
// Context on Cursor, where postToolUseFailure has no documented outputs.
type ignoredPostToolFailureResults struct{}

// Context returns a no-op result; Cursor ignores postToolUseFailure outputs.
func (ignoredPostToolFailureResults) Context(string) model.PostToolFailureResult {
	return ignoredPostToolFailureResult{}
}

type ignoredPostToolFailureResult struct{}

// IsZero reports that the discarded Cursor result carries no instruction.
func (ignoredPostToolFailureResult) IsZero() bool { return true }
