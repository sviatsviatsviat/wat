package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionRequest is the permissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name (VS Code).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the native tool input JSON (VS Code).
	ToolInput json.RawMessage `json:"tool_input"`
	// ToolArgs is the native tool input JSON (camelCase).
	ToolArgs json.RawMessage `json:"toolArgs"`
}

// EventName returns the canonical hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

// NativeToolName returns the tool name from either wire format.
func (e PermissionRequest) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input JSON from either wire format.
func (e PermissionRequest) Input() json.RawMessage {
	if len(e.ToolInput) > 0 {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e PermissionRequest) ShellCommand() string {
	if !isShellToolName(e.NativeToolName()) {
		return ""
	}
	return hookkit.ExtractShellCommand(e.Input())
}

// PermissionRequestOutput is the response for permissionRequest events.
type PermissionRequestOutput struct {
	// Behavior is allow or deny.
	Behavior string
	// Message is the permission message.
	Message string
	// Interrupt stops the session when true.
	Interrupt bool
	// SuppressWarnExit skips exit code 2 when Behavior is deny. Use for ask-style
	// responses that emit deny on the wire without warn-exit semantics.
	SuppressWarnExit bool
}

func (o PermissionRequestOutput) isZero() bool {
	return o.Behavior == "" && o.Message == "" && !o.Interrupt
}

// PermissionRequestResults is the hook-scoped response builder supplied to Chain handlers by registration.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
	// Ask returns an ask-style deny that suppresses warn-exit semantics.
	Ask(message string) PermissionRequestOutput
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

// Ask returns an ask-style deny that suppresses warn-exit semantics.
func (permissionRequestResults) Ask(message string) PermissionRequestOutput {
	return PermissionRequestOutput{Behavior: "deny", Message: message, SuppressWarnExit: true}
}

func encodePermissionRequest(o PermissionRequestOutput) ([]byte, int, error) {
	if o.Behavior == "" && o.Message == "" && !o.Interrupt {
		return nil, 0, nil
	}
	out := map[string]any{}
	if o.Behavior != "" {
		out["behavior"] = o.Behavior
	}
	if o.Message != "" {
		out["message"] = o.Message
	}
	if o.Interrupt {
		out["interrupt"] = true
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.Behavior == "deny" && !o.SuppressWarnExit {
		exitCode = WarnExit
	}
	return b, exitCode, err
}

func init() {
	registerDecoder(EventPermissionRequest, decodeAs[PermissionRequest])
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
