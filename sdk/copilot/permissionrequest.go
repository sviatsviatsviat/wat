package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionRequest is the PermissionRequest hook event.
type PermissionRequest struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input.
	ToolInput tools.Input `json:"-"`
}

// EventName returns the canonical hook event name.
func (PermissionRequest) EventName() string { return EventPermissionRequest }

// NativeToolName returns the tool name.
func (e PermissionRequest) NativeToolName() string {
	return e.ToolName
}

// Input returns tool input.
func (e PermissionRequest) Input() tools.Input {
	return e.ToolInput
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e PermissionRequest) ShellCommand() string {
	if !isShellToolName(e.NativeToolName()) {
		return ""
	}
	return hookkit.ExtractShellCommand(e.Input().Raw())
}

// PermissionRequestOutput is the response for PermissionRequest events.
// Construct via PermissionRequestResults builders and With* methods. A nil value is a no-op.
type PermissionRequestOutput interface {
	isPermissionRequestOutput()
	// WithInterrupt stops the session when true.
	WithInterrupt(v bool) PermissionRequestOutput
}

type permissionRequestOutput struct {
	behavior         string
	message          string
	interrupt        bool
	suppressWarnExit bool
}

func (permissionRequestOutput) isPermissionRequestOutput() {}

func (o permissionRequestOutput) isZero() bool {
	return o.behavior == "" && o.message == "" && !o.interrupt
}

// WithInterrupt stops the session when true.
func (o permissionRequestOutput) WithInterrupt(v bool) PermissionRequestOutput {
	o.interrupt = v
	return o
}

// PermissionRequestResults is the hook-scoped response builder supplied to On* handlers by registration.
type PermissionRequestResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionRequestOutput
	// Deny returns a deny verdict with a permission message.
	Deny(message string) PermissionRequestOutput
	// Ask returns an ask-style deny that suppresses warn-exit semantics.
	Ask(message string) PermissionRequestOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PermissionRequestOutput
	isPermissionRequestResults()
}

type permissionRequestResults struct{}

func (permissionRequestResults) isPermissionRequestResults() {}

// Allow returns an allow verdict.
func (permissionRequestResults) Allow() PermissionRequestOutput {
	return permissionRequestOutput{behavior: "allow"}
}

// Deny returns a deny verdict with a permission message.
func (permissionRequestResults) Deny(message string) PermissionRequestOutput {
	return permissionRequestOutput{behavior: "deny", message: message}
}

// Ask returns an ask-style deny that suppresses warn-exit semantics.
func (permissionRequestResults) Ask(message string) PermissionRequestOutput {
	return permissionRequestOutput{behavior: "deny", message: message, suppressWarnExit: true}
}

// Noop returns an empty response (silent stdout).
func (permissionRequestResults) Noop() PermissionRequestOutput {
	return permissionRequestOutput{}
}

func (permissionRequestOutput) allowedEvents() []string {
	return []string{EventPermissionRequest}
}

func (o permissionRequestOutput) encode() ([]byte, int, error) {
	if o.behavior == "" && o.message == "" && !o.interrupt {
		return nil, 0, nil
	}
	out := map[string]any{}
	if o.behavior != "" {
		out["behavior"] = o.behavior
	}
	if o.message != "" {
		out["message"] = o.message
	}
	if o.interrupt {
		out["interrupt"] = true
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.behavior == "deny" && !o.suppressWarnExit {
		exitCode = WarnExit
	}
	return b, exitCode, err
}

func init() {
	registerDecoder(EventPermissionRequest, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PermissionRequest, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.NativeToolName(), raw, "tool_input")
		})
	})
}

// OnPermissionRequest registers a PermissionRequest handler.
func OnPermissionRequest(fn func(context.Context, Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	return (&chain{}).PermissionRequest(fn)
}

// PermissionRequest registers another PermissionRequest handler on the chain.
func (c *chain) PermissionRequest(fn func(context.Context, Hook[PermissionRequest], PermissionRequestResults) (PermissionRequestOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PermissionRequest) (PermissionRequestOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), permissionRequestResults{})
	})
	return c
}
