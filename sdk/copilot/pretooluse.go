package copilot

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision string

const (
	// DecisionAllow permits the tool call.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the tool call.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
)

// PreToolUse is the preToolUse hook event.
type PreToolUse struct {
	Envelope
	// ToolName is the tool name (VS Code snake_case).
	ToolName string `json:"tool_name"`
	// ToolNameCamel is the tool name (camelCase).
	ToolNameCamel string `json:"toolName"`
	// ToolInput is the typed tool input (VS Code).
	ToolInput tools.Input `json:"-"`
	// ToolArgs is the typed tool input (camelCase).
	ToolArgs tools.Input `json:"-"`
}

// EventName returns the canonical hook event name.
func (PreToolUse) EventName() string { return EventPreToolUse }

// NativeToolName returns the tool name from either wire format.
func (e PreToolUse) NativeToolName() string {
	if e.ToolName != "" {
		return e.ToolName
	}
	return e.ToolNameCamel
}

// Input returns tool input from either wire format.
func (e PreToolUse) Input() tools.Input {
	if e.ToolInput.HasRaw() {
		return e.ToolInput
	}
	return e.ToolArgs
}

// ShellCommand extracts the shell command when the tool is a shell execution tool.
func (e PreToolUse) ShellCommand() string {
	if !isShellToolName(e.NativeToolName()) {
		return ""
	}
	return hookkit.ExtractShellCommand(e.Input().Raw())
}

// PreToolOutput is the response for preToolUse events.
type PreToolOutput struct {
	// Decision is the permission verdict (allow, deny, ask).
	Decision PermissionDecision
	// Reason is the agent-facing decision reason.
	Reason string
	// ModifiedArgs replaces tool arguments when set.
	ModifiedArgs map[string]any
}

func (o PreToolOutput) isZero() bool {
	return o.Decision == "" && o.Reason == "" && o.ModifiedArgs == nil
}

// PreToolResults is the hook-scoped response builder supplied to Chain.PreToolUse handlers by registration.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolOutput
	isPreToolResults()
}

type preToolResults struct{}

func (preToolResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (preToolResults) Allow() PreToolOutput {
	return PreToolOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolResults) Deny(reason string) PreToolOutput {
	return PreToolOutput{Decision: DecisionDeny, Reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolResults) Ask(reason string) PreToolOutput {
	return PreToolOutput{Decision: DecisionAsk, Reason: reason}
}

func encodePreTool(o PreToolOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Decision != "" {
		out["permissionDecision"] = string(o.Decision)
		if o.Reason != "" {
			out["permissionDecisionReason"] = o.Reason
		}
	}
	if o.ModifiedArgs != nil {
		out["modifiedArgs"] = o.ModifiedArgs
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	return b, 0, err
}

func init() {
	registerDecoder(EventPreToolUse, func(raw []byte, received, canonical string) (Event, error) {
		return decodeAsAndThen(raw, received, canonical, func(e *PreToolUse, raw []byte) {
			name := e.NativeToolName()
			e.ToolInput = tools.NewInputFromPayload(name, raw, "tool_input")
			e.ToolArgs = tools.NewInputFromPayload(name, raw, "toolArgs")
		})
	})
}

// PreToolUse registers a PreToolUse handler.
func (c *Chain) PreToolUse(fn func(context.Context, PreToolUseHook, PreToolResults) (PreToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PreToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preToolResults{})
	})
	return &Chain{}
}
