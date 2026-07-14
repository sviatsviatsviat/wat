package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
}

// EventName returns the hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

func init() {
	registerDecoder(EventPermissionRequest, decodeAs[PermissionRequest])
}

// PermissionRequestOutput is the response for PermissionRequest events.
type PermissionRequestOutput struct {
	Common
	// Behavior is allow or deny.
	Behavior string
	// UpdatedInput replaces tool arguments when set.
	UpdatedInput map[string]any
	// Message is the permission message.
	Message string
	// Interrupt stops the session when true.
	Interrupt bool
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o PermissionRequestOutput) isZero() bool {
	return o.Common.isZero() && o.Behavior == "" && o.UpdatedInput == nil &&
		o.Message == "" && !o.Interrupt && o.AdditionalContext == ""
}

// PermissionRequestResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
	isPermissionRequestResults()
}

type permissionRequestResults struct{}

func (permissionRequestResults) isPermissionRequestResults() {}

// Allow returns an allow verdict.
func (permissionRequestResults) Allow() PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "allow"}
}

// Deny returns a deny verdict with a permission message.
func (permissionRequestResults) Deny(message string) PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "deny", Message: message}
}

// PermissionRequest registers a PermissionRequest handler.
func (c *Chain) PermissionRequest(fn func(context.Context, PermissionRequestHook, PermissionRequestResults) (PermissionRequestOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionRequest) (PermissionRequestOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionRequestResults{})
	})
	return &Chain{}
}
