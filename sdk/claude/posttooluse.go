package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/claude/tools"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PostToolUse is the PostToolUse hook event.
type PostToolUse struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the typed tool input for ToolName.
	ToolInput tools.Input `json:"-"`
	// ToolUseID is the tool use identifier.
	ToolUseID string `json:"tool_use_id"`
	// ToolResponse is the tool response JSON.
	ToolResponse json.RawMessage `json:"tool_response"`
	// DurationMs is the tool execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// EventName returns the hook event name.
func (PostToolUse) EventName() string { return EventPostToolUse }

func init() {
	codec.Register(EventPostToolUse, func(raw []byte) (run.Event, error) {
		return hookkit.DecodeEvent(codec, raw, func(e *PostToolUse, raw []byte) {
			e.ToolInput = tools.NewInputFromPayload(e.ToolName, raw, "tool_input")
		})
	})
}

// PostToolUseOutput is the response for PostToolUse and PostToolUseFailure events.
// Construct via PostToolUseResults / PostToolUseFailureResults builders and With* methods.
// A nil value is a no-op.
type PostToolUseOutput interface {
	run.Output
	isPostToolUseOutput()
	// WithUpdatedToolOutput replaces the tool result when set.
	WithUpdatedToolOutput(output any) PostToolUseOutput
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) PostToolUseOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) PostToolUseOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) PostToolUseOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) PostToolUseOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) PostToolUseOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) PostToolUseOutput
}

type postToolUseOutput struct {
	common
	eventName         string
	block             bool
	reason            string
	additionalContext string
	updatedToolOutput any
}

func (postToolUseOutput) isPostToolUseOutput() {}

// IsZero reports whether this hook response is empty.
func (o postToolUseOutput) IsZero() bool {
	return o.common.IsZero() && !o.block && o.reason == "" &&
		o.additionalContext == "" && o.updatedToolOutput == nil
}

// WithUpdatedToolOutput replaces the tool result when set.
func (o postToolUseOutput) WithUpdatedToolOutput(output any) PostToolUseOutput {
	o.updatedToolOutput = hookkit.CloneJSONValue(output)
	return o
}

// WithAdditionalContext injects model context.
func (o postToolUseOutput) WithAdditionalContext(text string) PostToolUseOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o postToolUseOutput) WithContinue(v bool) PostToolUseOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o postToolUseOutput) WithStopReason(reason string) PostToolUseOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o postToolUseOutput) WithSuppressOutput(v bool) PostToolUseOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o postToolUseOutput) WithSystemMessage(msg string) PostToolUseOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o postToolUseOutput) WithTerminalSequence(seq string) PostToolUseOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// PostToolUseResults is the hook-scoped response builder supplied to On* handlers by registration.
type PostToolUseResults interface {
	// Context returns a context-injection-only PostToolUse result.
	Context(text string) PostToolUseOutput
	// Block returns a block result with an agent-facing reason.
	Block(reason string) PostToolUseOutput
	isPostToolUseResults()
}

type postToolUseResults struct{}

func (postToolUseResults) isPostToolUseResults() {}

// Context returns a context-injection-only PostToolUse result.
func (postToolUseResults) Context(text string) PostToolUseOutput {
	return postToolUseOutput{eventName: EventPostToolUse, additionalContext: text}
}

// Block returns a block result with an agent-facing reason.
func (postToolUseResults) Block(reason string) PostToolUseOutput {
	return postToolUseOutput{eventName: EventPostToolUse, block: true, reason: reason}
}

func (o postToolUseOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.updatedToolOutput != nil {
		hso["updatedToolOutput"] = o.updatedToolOutput
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o postToolUseOutput) Encode() ([]byte, int, error) {
	return marshalHookOutput(o.eventName, o.encodeInto)
}

// Merge combines other into this PostToolUse output.
func (o postToolUseOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(postToolUseOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.common.Merge(b.common)
	oDec, bDec := "", ""
	if o.block {
		oDec = "block"
	}
	if b.block {
		bDec = "block"
	}
	dec, reason := hookkit.MergeRankedString(oDec, o.reason, bDec, b.reason, hookkit.BlockDecisionRankString)
	updatedToolOutput, w := hookkit.TakeLastAny("updatedToolOutput", o.updatedToolOutput, b.updatedToolOutput)
	if w != "" {
		warnings = append(warnings, w)
	}
	eventName := o.eventName
	if eventName == "" {
		eventName = b.eventName
	}
	return postToolUseOutput{
		common:            mergedCommon,
		eventName:         eventName,
		block:             dec == "block",
		reason:            reason,
		additionalContext: hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
		updatedToolOutput: updatedToolOutput,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o postToolUseOutput) Stop() bool {
	return o.common.Stop() || o.block
}

// PostToolUse registers a PostToolUse handler on the chain.
func (c *chain) PostToolUse(fn func(context.Context, run.Hook[PostToolUse], PostToolUseResults) (PostToolUseOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[PostToolUse]) (PostToolUseOutput, error) {
		return fn(ctx, hook, postToolUseResults{})
	}))
	return c
}
