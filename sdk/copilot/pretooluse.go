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
// Construct via PreToolResults builders and With* methods. A nil value is a no-op.
type PreToolOutput interface {
	isPreToolOutput()
	// WithModifiedArgs replaces tool arguments when set.
	WithModifiedArgs(args map[string]any) PreToolOutput
}

type preToolOutput struct {
	decision     PermissionDecision
	reason       string
	modifiedArgs map[string]any
}

func (preToolOutput) isPreToolOutput() {}

func (o preToolOutput) isZero() bool {
	return o.decision == "" && o.reason == "" && o.modifiedArgs == nil
}

// WithModifiedArgs replaces tool arguments when set.
func (o preToolOutput) WithModifiedArgs(args map[string]any) PreToolOutput {
	o.modifiedArgs = args
	return o
}

// PreToolResults is the hook-scoped response builder supplied to Chain.PreToolUse handlers by registration.
type PreToolResults interface {
	// Allow returns an allow verdict.
	Allow() PreToolOutput
	// Deny returns a deny verdict with an agent-facing reason.
	Deny(reason string) PreToolOutput
	// Ask returns an ask verdict with an agent-facing reason.
	Ask(reason string) PreToolOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PreToolOutput
	isPreToolResults()
}

type preToolResults struct{}

func (preToolResults) isPreToolResults() {}

// Allow returns an allow verdict.
func (preToolResults) Allow() PreToolOutput {
	return preToolOutput{decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing reason.
func (preToolResults) Deny(reason string) PreToolOutput {
	return preToolOutput{decision: DecisionDeny, reason: reason}
}

// Ask returns an ask verdict with an agent-facing reason.
func (preToolResults) Ask(reason string) PreToolOutput {
	return preToolOutput{decision: DecisionAsk, reason: reason}
}

// Noop returns an empty response (silent stdout).
func (preToolResults) Noop() PreToolOutput {
	return preToolOutput{}
}

func (preToolOutput) allowedEvents() []string {
	return []string{EventPreToolUse}
}

func (o preToolOutput) encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.decision != "" {
		out["permissionDecision"] = string(o.decision)
		if o.reason != "" {
			out["permissionDecisionReason"] = o.reason
		}
	}
	if o.modifiedArgs != nil {
		out["modifiedArgs"] = o.modifiedArgs
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
func (c *Chain) PreToolUse(fn func(context.Context, Hook[PreToolUse], PreToolResults) (PreToolOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev PreToolUse) (PreToolOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), preToolResults{})
	})
	return &Chain{}
}
