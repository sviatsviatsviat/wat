package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// MessageDisplay is the MessageDisplay hook event.
type MessageDisplay struct {
	Envelope
	// TurnID is the turn identifier.
	TurnID string `json:"turn_id"`
	// MessageID is the message identifier.
	MessageID string `json:"message_id"`
	// Index is the message index in the turn.
	Index int `json:"index"`
	// Final is true when this is the final delta.
	Final bool `json:"final"`
	// Delta is the streamed message delta.
	Delta string `json:"delta"`
}

// EventName returns the hook event name.
func (MessageDisplay) EventName() string { return EventMessageDisplay }

func init() {
	codec.Register(EventMessageDisplay, hookkit.EventDecoder[MessageDisplay](codec))
}

// MessageDisplayOutput is the response for MessageDisplay events.
// Construct via MessageDisplayResults builders and With* methods.
// A nil value is a no-op.
type MessageDisplayOutput interface {
	run.Output
	isMessageDisplayOutput()
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) MessageDisplayOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) MessageDisplayOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) MessageDisplayOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) MessageDisplayOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) MessageDisplayOutput
}

type messageDisplayOutput struct {
	common
	displayContent *string
}

func (messageDisplayOutput) isMessageDisplayOutput() {}

// IsZero reports whether this hook response is empty.
func (o messageDisplayOutput) IsZero() bool {
	return o.common.IsZero() && o.displayContent == nil
}

// WithContinue sets whether Claude should continue the session.
func (o messageDisplayOutput) WithContinue(v bool) MessageDisplayOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o messageDisplayOutput) WithStopReason(reason string) MessageDisplayOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o messageDisplayOutput) WithSuppressOutput(v bool) MessageDisplayOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o messageDisplayOutput) WithSystemMessage(msg string) MessageDisplayOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o messageDisplayOutput) WithTerminalSequence(seq string) MessageDisplayOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// MessageDisplayResults is the hook-scoped response builder supplied to On* handlers by registration.
type MessageDisplayResults interface {
	// Override returns a display-content override result.
	Override(content string) MessageDisplayOutput
	isMessageDisplayResults()
}

type messageDisplayResults struct{}

func (messageDisplayResults) isMessageDisplayResults() {}

// Override returns a display-content override result.
func (messageDisplayResults) Override(content string) MessageDisplayOutput {
	c := content
	return messageDisplayOutput{displayContent: &c}
}

func (o messageDisplayOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.displayContent != nil {
		hso["displayContent"] = *o.displayContent
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o messageDisplayOutput) Encode() ([]byte, int, error) {
	return marshalHookOutput(EventMessageDisplay, o.encodeInto)
}

// Merge combines other into this MessageDisplay output.
func (o messageDisplayOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(messageDisplayOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.common.Merge(b.common)
	displayContent, w := hookkit.TakeLastPtr("displayContent", o.displayContent, b.displayContent)
	if w != "" {
		warnings = append(warnings, w)
	}
	return messageDisplayOutput{
		common:         mergedCommon,
		displayContent: displayContent,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o messageDisplayOutput) Stop() bool {
	return o.common.Stop()
}

// MessageDisplay registers a MessageDisplay handler on the chain.
func (c *chain) MessageDisplay(fn func(context.Context, run.Hook[MessageDisplay], MessageDisplayResults) (MessageDisplayOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[MessageDisplay]) (MessageDisplayOutput, error) {
		return fn(ctx, hook, messageDisplayResults{})
	}))
	return c
}
