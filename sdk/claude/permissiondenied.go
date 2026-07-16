package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionDenied is the PermissionDenied hook event.
type PermissionDenied struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// Reason is the classifier denial reason.
	Reason string `json:"reason"`
}

// EventName returns the hook event name.
func (PermissionDenied) EventName() string { return EventPermissionDenied }

func init() {
	registerDecoder(EventPermissionDenied, decodeAs[PermissionDenied])
}

// PermissionDeniedOutput is the response for PermissionDenied events.
type PermissionDeniedOutput struct {
	Common
	// Retry requests a permission retry when true.
	Retry bool
}

func (o PermissionDeniedOutput) isZero() bool {
	return o.Common.isZero() && !o.Retry
}

// PermissionDenied registers a PermissionDenied handler.
func (c *Chain) PermissionDenied(fn func(context.Context, PermissionDeniedHook, PermissionDeniedResults) (PermissionDeniedOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionDenied) (PermissionDeniedOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionDeniedResults{})
	})
	return &Chain{}
}

// PermissionDeniedResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PermissionDeniedResults interface {
	// Retry returns a retry-requested PermissionDenied result.
	Retry() PermissionDeniedOutput
	isPermissionDeniedResults()
}

type permissionDeniedResults struct{}

func (permissionDeniedResults) isPermissionDeniedResults() {}

// Retry returns a retry-requested PermissionDenied result.
func (permissionDeniedResults) Retry() PermissionDeniedOutput {
	return PermissionDeniedOutput{Retry: true}
}
