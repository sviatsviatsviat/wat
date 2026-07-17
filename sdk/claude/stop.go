package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Stop is the Stop hook event.
type Stop struct {
	Envelope
	// StopHookActive is true when a stop hook already forced continuation.
	StopHookActive bool `json:"stop_hook_active"`
	// LastAssistantMessage is the final assistant text of the turn.
	LastAssistantMessage string `json:"last_assistant_message"`
}

// EventName returns the hook event name.
func (Stop) EventName() string { return EventStop }

func init() {
	registerDecoder(EventStop, decodeAs[Stop])
}

// StopOutput is the response for Stop and SubagentStop events.
// Construct via StopResults builders and With* methods.
// A nil value is a no-op.
type StopOutput interface {
	isStopOutput()
	// WithAdditionalContext is non-error feedback that continues the conversation.
	WithAdditionalContext(text string) StopOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) StopOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) StopOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) StopOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) StopOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) StopOutput
}

type stopOutput struct {
	common
	block             bool
	reason            string
	additionalContext string
}

func (stopOutput) isStopOutput() {}
func (o stopOutput) isZero() bool {
	return o.common.isZero() && !o.block && o.reason == "" && o.additionalContext == ""
}

// WithAdditionalContext is non-error feedback that continues the conversation.
func (o stopOutput) WithAdditionalContext(text string) StopOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o stopOutput) WithContinue(v bool) StopOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o stopOutput) WithStopReason(reason string) StopOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o stopOutput) WithSuppressOutput(v bool) StopOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o stopOutput) WithSystemMessage(msg string) StopOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o stopOutput) WithTerminalSequence(seq string) StopOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// StopResults is the hook-scoped response builder supplied to Chain handlers by registration.
type StopResults interface {
	// Context returns non-blocking feedback that continues the conversation.
	Context(text string) StopOutput
	// FollowUp blocks completion and feeds reason back to Claude.
	FollowUp(reason string) StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// Context returns non-blocking feedback that continues the conversation.
func (stopResults) Context(text string) StopOutput {
	return stopOutput{additionalContext: text}
}

// FollowUp blocks completion and feeds reason back to Claude.
func (stopResults) FollowUp(reason string) StopOutput {
	return stopOutput{block: true, reason: reason}
}

func (stopOutput) allowedEvents() []string {
	return []string{EventStop, EventSubagentStop}
}

func (o stopOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.block {
		top["decision"] = "block"
		if o.reason != "" {
			top["reason"] = o.reason
		}
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Stop registers a Stop handler.
func (c *Chain) Stop(fn func(context.Context, Hook[Stop], StopResults) (StopOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Stop) (StopOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), stopResults{})
	})
	return &Chain{}
}
